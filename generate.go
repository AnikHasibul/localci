package localci

import (
	"encoding/json"
	"log"
	"os"
)

// Generate generates a config file
func Generate() {
	c := make(config)
	c["0Build"] = []execs{
		execs{
			"GoBuild": []string{"go", "build"},
		},
	}
	c["1Test"] = []execs{
		execs{
			"GoTest": []string{"go", "test"},
		},
	}
	c["2Run"] = []execs{
		execs{
			"GoRun": []string{"go", "run", "."},
		},
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}
