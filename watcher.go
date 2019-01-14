package localci

import (
	"os"
	"time"
)

// watch calls the notifier after every 1 second
func (ci *ciObj) watch() {
	for {
		<-time.After(1 * time.Second)
		ci.notify()
	}
}

// addToWatcher adds given files to the watcher map
func (ci *ciObj) addToWatcher() {
	for _, v := range ci.files {
		var fi os.FileInfo
		fi, ci.err = os.Stat(v)
		if ci.err != nil {
			continue
		}
		ci.fs[v] = fi.ModTime()
	}
}

// waitForEvent returns the name of changed file
// it doesn't return until a file changes
func (ci *ciObj) waitForEvent() string {
	return <-ci.fileTicker
}

// listen listens for an event
// calls jobs, sets queue and...
func (ci *ciObj) listen() {
	for {
		e := ci.waitForEvent()
		ci.queue.Add()
		ci.sessionLog()
		ci.vLog(e)
		go ci.jobs()
	}
}
