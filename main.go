package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	githubactions "github.com/sethvargo/go-githubactions"
)

var (
	// errGoNotFound where is go??
	errGoNotFound = errors.New("could not find go in PATH")

	// errRunFailed means we could not run the thing
	errRunFailed = errors.New("failed to execute go test -cover -json")
)

type jsonfile []string

func (j *jsonfile) Write(p []byte) (n int, err error) {
	var l = strings.Split(string(p), "\n")

	for i := range l {
		if l[i] == "" {
			continue
		}

		*j = append(*j, l[i])
	}

	return len(p), nil
}

func (j jsonfile) String() string {
	return "[" + strings.Join(j, ",") + "]"
}

func run() (string, error) {
	var (
		err error
		g   string
		c   *exec.Cmd
		j   jsonfile
	)

	g, err = exec.LookPath("go")
	if err != nil {
		return "", errGoNotFound
	}

	c = exec.Command(g, "test", "-cover", "-json", "./...")
	c.Stdout = &j

	err = c.Run()
	if err != nil {
		return "", fmt.Errorf("%w: %s", errRunFailed, err.Error())
	}

	return j.String(), nil
}

func main() {
	var j, err = run()
	if err != nil {
		log.Println(err.Error())
		githubactions.Fatalf(err.Error())
	}

	githubactions.SetOutput("gotessa_json", j)
}
