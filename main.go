package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	shellquote "github.com/kballard/go-shellquote"
)

type config struct {
	AdditionalParams string `env:"additional_params"`
	ProjectLocation  string `env:"project_location,dir"`
}

func failf(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
	os.Exit(1)
}

func hasAnalyzeError(cmdOutput string) bool {
	// example: error • Undefined class 'function' • lib/package.dart:3:1 • undefined_class
	analyzeErrorPattern := regexp.MustCompile(`error.+\.dart:\d+:\d+`)
	if analyzeErrorPattern.MatchString(cmdOutput) {
		return true
	}

	return false
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(cfg)

	additionalParams, err := shellquote.Split(cfg.AdditionalParams)
	if err != nil {
		failf("Failed to parse additional parameters, error: %s", err)
	}

	fmt.Println()
	log.Infof("Running analyze")

	var b bytes.Buffer
	multiwr := io.MultiWriter(os.Stdout, &b)
	analyzeCmd := command.New("flutter", append([]string{"analyze"}, additionalParams...)...).
		SetDir(cfg.ProjectLocation).
		SetStdout(multiwr).
		SetStderr(os.Stderr)

	fmt.Println()
	log.Donef("$ %s", analyzeCmd.PrintableCommandArgs())
	fmt.Println()
	
	if err := analyzeCmd.Run(); err != nil {
		if hasAnalyzeError(b.String()) {
			log.Errorf("flutter analyze found errors: %s", err)
			os.Exit(1)
		}
	}
}
