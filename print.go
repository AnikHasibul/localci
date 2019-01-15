package localci

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
)

// print the pass text
func (ci *ciObj) success() {
	fmt.Print(aurora.Green("PASS"), "\n")
}

// print the fail text
func (ci *ciObj) fail() {
	fmt.Print(aurora.Red("FAIL"), "\n")
}

// print the cancel text
func (ci *ciObj) cancel() {
	fmt.Print(aurora.Red("CANCEL"), "\n")
}

// print a normal message
func (ci *ciObj) msg(msg string) {
	if ci.verbose {
		fmt.Println(msg)
	}
}

// print a verbose log
func (ci *ciObj) vLog(msg string) {
	if ci.verbose {
		fmt.Println(msg)
	}
}

// print a log
func (ci *ciObj) log(msg string) {
	fmt.Println(msg)
}

// print the stage log
func (ci *ciObj) stageLog(stg string) {
	fmt.Printf(
		"%v Running...\n",
		aurora.Blue(stg),
	)
}

// print a job log
func (ci *ciObj) jobLog(stg string) {
	fmt.Printf(
		"%v ...",
		aurora.Magenta(stg),
	)
}

// print a version log
func (ci *ciObj) sessionLog() {
	h, m, s := time.Now().Clock()
	fmt.Printf(
		"\n%v ",
		aurora.Brown(fmt.Sprintf("v%d:%d:%d", h, m, s)),
	)
}
