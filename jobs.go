package localci

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"

	yml "gopkg.in/yaml.v2"
)

// parseConfig parses the config file
func (ci *ciObj) parseConfig() {
	// read the config file
	c, err := ioutil.ReadFile(ci.configFile)
	if err != nil {
		ci.err = err
		return
	}
	// fix #1
	// delete the previous config
	ci.config = config{}
	// parse it
	ci.err = yml.Unmarshal(c, &ci.config)
	if err != nil {
		return
	}
}

// runStages runs all stages
func (ci *ciObj) runStages() {
	// loop over all stages
	for _, execs := range ci.config {
		// loop over all jobs on the stages
		for stage, e := range execs {
			ci.stageLog(stage)
			for _, v := range e {
				ci.jobLog(v)
				// run each job
				ci.run(v)
				if ci.err == nil {
					ci.success()
				} else {
					ci.fail()
					ci.msg(ci.err.Error())
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
func (ci *ciObj) run(cmdLine string) {
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
		"sh",
		"-c",
		cmdLine,
	)
	ci.cmd.Stdout = os.Stdout
	ci.cmd.Stderr = os.Stderr
	ci.err = ci.cmd.Run()

	if ci.err != nil {
		ci.err = fmt.Errorf("%s: %v", cmdLine, ci.err)
		return
	}

	if !ci.cmd.ProcessState.Success() {
		ci.err = fmt.Errorf("%s: %s", cmdLine, "exited with !0")
	}
}

// jobs calls all the jobs needed to do
func (ci *ciObj) jobs() {
	// release queue
	defer ci.queue.Done()
	// run....
	ci.runStages()
	if ci.err != nil {
		ci.log(ci.err.Error())
	}
}
