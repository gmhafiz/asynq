package cli

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/rogpeppe/go-internal/semver"
	"github.com/spf13/cobra"
)

// https://gist.github.com/ik5/d8ecde700972d4378d87#file-colors-go
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var rootCmd = &cobra.Command{
	Use:   "asynqgen",
	Short: "asynqgen is an asynq task scaffolder",
	Long:  `Quickly scaffold a new task - https://github.com/gmhafiz/asynq`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func checkVersion() error {
	version := runtime.Version()

	goSupport := semver.Compare(version, "1.16")
	if goSupport < 0 {
		return errors.New("warning: only last 2 versions of go are supported officially by the Go team. See https://golang.org/doc/devel/release#policy")
	}

	return nil
}
