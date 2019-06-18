package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
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

	analyzeCmd := command.New("flutter", append([]string{"analyze"}, additionalParams...)...).
		SetDir(cfg.ProjectLocation)

	fmt.Println()
	log.Donef("$ %s", analyzeCmd.PrintableCommandArgs())
	fmt.Println()

	out, err := analyzeCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			// example: error • Undefined class 'function' • lib/package.dart:3:1 • undefined_class
			analyzeErrorPattern := regexp.MustCompile(`error.+\.dart.\s*`)
			if !analyzeErrorPattern.MatchString(out) {
				// false positive, no error in analyze output
				// flutter analyze returns with nonzero for 'info'
				// level errors, see: https://github.com/flutter/flutter/issues/20855
				log.Printf(out)
				os.Exit(0)
			}
		}

		log.Printf(out)
		failf("Running command failed, error: %s", err)
	}

	log.Printf(out)
}
