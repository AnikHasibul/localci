package localci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
)

// parseConfig parses the config file
func (ci *ciObj) parseConfig() {
	// read the config file
	c, err := ioutil.ReadFile(ci.configFile)
	if err != nil {
		ci.err = err
		return
	}
	// parse it
	ci.err = json.Unmarshal(c, &ci.config)
	if err != nil {
		return
	}
}

// runStages runs all stages
func (ci *ciObj) runStages() {
	// loop over all stages
	for stage, execs := range ci.config {
		ci.stageLog(stage)
		// loop over all jobs on the stages
		for _, e := range execs {
			for k, v := range e {
				ci.jobLog(k)
				// run each job
				ci.run(v)
				if ci.err == nil {
					ci.success()
				} else {
					ci.fail()
					ci.msg(ci.err.Error())
				}
				if ci.writeToStdout {
					if w := string(ci.flush()); w != "" {
						ci.log(w)
					}
				}
				if ci.err != nil {
					ci.err = nil
					return
				}
			}
		}
	}
}

// run runs a command
func (ci *ciObj) run(args []string) {
	// if a process exists on this pid
	// kill it
	if ci.cmd != nil {
		if ci.cmd.Process != nil {
			// nolint
			ci.cmd.Process.Kill()
		}
	}
	var done = make(chan bool, 1)
	defer close(done)
	go func() {
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt)
		defer signal.Stop(kill)
		select {
		case <-kill:
			if ci.cmd != nil {
				if ci.cmd.Process != nil {
					ci.cancel()
					// nolint
					ci.cmd.Process.Kill()
				}
			}
		case <-done:
			return
		}
	}()
	// new command
	ci.cmd = exec.Command(
		args[0],
		args[1:]...,
	)
	ci.cmd.Stdout = ci
	ci.cmd.Stderr = ci
	ci.err = ci.cmd.Run()

	if ci.err != nil {
		ci.err = fmt.Errorf("%s: %v", args[0], ci.err)
		return
	}

	if !ci.cmd.ProcessState.Success() {
		ci.err = fmt.Errorf("%s: %s", args[0], "exited with !0")
	}
}

// jobs calls all the jobs needed to do
func (ci *ciObj) jobs() {
	// release queue
	defer ci.queue.Done()
	// parse config
	ci.parseConfig()
	if ci.err != nil {
		ci.log(ci.err.Error())
	}
	// run....
	ci.runStages()
	if ci.err != nil {
		ci.log(ci.err.Error())
	}
}
