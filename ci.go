package localci

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/anikhasibul/queue"
)

// ciObj holds the whole ci methods
type ciObj struct {
	config        config
	cmd           *exec.Cmd
	queue         *queue.Q
	fs            watchMap
	err           error
	fileTicker    chan string
	configFile    string
	files         []string
	wBuff         []byte
	verbose       bool
	writeToStdout bool
}

type (
	// config format
	// map[stages][]execs
	config map[string][]execs
	// exec format
	// map[title]args
	execs map[string][]string
)

// initialize initializes the whole ci process
func initialize() *ciObj {
	ci := new(ciObj)
	//
	flag.BoolVar(
		&ci.writeToStdout,
		"stdout",
		true,
		"Write linters output to stdout.",
	)
	flag.BoolVar(
		&ci.verbose,
		"v",
		true,
		"Verbose output.",
	)
	flag.StringVar(
		&ci.configFile,
		"c",
		".lci.json",
		"Config file.",
	)
	gen := flag.Bool(
		"gen",
		false,
		"generates a config file. To save the config, try: $ localci -gen >.lci.json",
	)
	flag.Parse()
	// get filea
	if *gen {
		Generate()
		os.Exit(0)
	}
	ci.files = flag.Args()
	// queue group (maximum 1 ci)
	ci.queue = queue.New(1)
	ci.fs = make(watchMap)
	ci.fileTicker = make(chan string, 50)
	return ci
}

// Start starts the ci
func Start() {
	// initialize the ci
	ci := initialize()
	if ci.err != nil {
		log.Fatal(ci.err)
	}
	// add files to watcher
	ci.addToWatcher()
	if ci.err != nil {
		log.Fatal(ci.err)
	}
	// start watching
	go ci.watch()
	// start listening to events
	ci.listen()
}
