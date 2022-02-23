package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/kballard/go-shellquote"
)

const (
	errorLevel   = "error"
	warningLevel = "warning"
	infoLevel    = "info"

	newLine = "\n"
)

var severityRegExp = map[string]string{
	errorLevel:   "error",
	warningLevel: "(error|warning)",
	infoLevel:    "(error|warning|info)",
}

type config struct {
	AdditionalParams string `env:"additional_params"`
	ProjectLocation  string `env:"project_location,dir"`
	FailSeverity     string `env:"fail_severity,opt[error,warning,info]"`
}

func failf(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
	os.Exit(1)
}

func constructRegex(severityPattern string) *regexp.Regexp {
	pattern := fmt.Sprintf(`^%s .+\.dart:\d+:\d+`, severityPattern)
	return regexp.MustCompile(pattern)
}

func hasAnalyzeError(cmdOutput string, failSeverity string) bool {
	// example: error • Undefined class 'function' • lib/package.dart:3:1 • undefined_class
	outputLines := strings.Split(cmdOutput, newLine)
	analyzeErrorPattern := constructRegex(severityRegExp[failSeverity])

	for _, line := range outputLines {
		if analyzeErrorPattern.MatchString(strings.TrimSpace(line)) {
			return true
		}
	}

	return false
}

func hasOtherError(cmdOutput string) bool {
	return !hasAnalyzeError(cmdOutput, infoLevel)
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Process config: failed to parse step inputs: %s", err)
	}
	stepconf.Print(cfg)

	additionalParams, err := shellquote.Split(cfg.AdditionalParams)
	if err != nil {
		failf("Process config: failed to parse additional parameters: %s", err)
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
		if hasAnalyzeError(b.String(), cfg.FailSeverity) {
			failf("Run: flutter analyze found errors: %s", err)
		} else if hasOtherError(b.String()) {
			failf("Run: step failed with error: %s", err)
		}
	}
}
