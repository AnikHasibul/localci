package localci

import (
	"log"
	"os"

	yml "gopkg.in/yaml.v2"
)

// Generate generates a config file
func Generate() {
	var c config
	c = append(c, map[string][]string{
		"test": []string{
			"go vet ./...",
			"go test ./...",
		},
	})
	c = append(c, map[string][]string{
		"build": []string{
			"go build ./...",
		},
	})
	b, err := yml.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}
