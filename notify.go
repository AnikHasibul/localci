package localci

import (
	"os"
	"time"
)

// watchMap holds the filename and last midified time
type watchMap map[string]time.Time

// notify sends an event if any file in the watchmap changed
func (ci *ciObj) notify() {
	for name, time := range ci.fs {
		var fi os.FileInfo
		// get the file info
		fi, ci.err = os.Stat(name)
		if ci.err != nil {
			continue
		}
		// if it's not changed
		// ignore
		if !fi.ModTime().After(time) {
			continue
		}
		// ignore directories
		if fi.IsDir() {
			continue
		}
		// update the map
		ci.fs[name] = fi.ModTime()
		// notify
		ci.fileTicker <- name
	}
}
