package console

import (
	"fmt"
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.kumparan.com/yowez/skeleton-service/config"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI
			applications in go using Cobra CLI library`,
}

// Execute :nodoc:
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.GetConf()
	setupLogger()
}

func setupLogger() {
	var formatter runtime.Formatter
	formatter = runtime.Formatter{
		ChildFormatter: &log.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if config.Env() == "development" {
		formatter = runtime.Formatter{
			ChildFormatter: &log.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)

	logLevel, err := log.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	hook, err := logrus_sentry.NewSentryHook(config.SentryDSN(), []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})
	if err != nil {
		log.Info("Logger configured to use only local stdout")
		return
	}

	hook.SetEnvironment(config.Env())
	hook.Timeout = 0 // fire and forget
	hook.StacktraceConfiguration.Enable = true
	log.AddHook(hook)
}
