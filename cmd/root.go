package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// flag for enabling verbose mode
	verbose bool
	// flag to show version and quit
	version bool
	// remote log sentry token
	sentryDSN string
	// BuildVersion to show in usage message
	BuildVersion string
	// prometheus metrics port
	metricsPort int

	// other flags
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "disco-bit",
	Short: "Lists commands",
	Long:  `Lists all available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			// only show version then exit
			fmt.Printf("%s build: %s", os.Args[0], BuildVersion)
			os.Exit(0)
		}
		fmt.Printf("\n\n")
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "default.flag.toml", "config file")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display debug messages")
	RootCmd.Flags().BoolVar(&version, "version", false, "Display app version and quit")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//aCmd.Flags().StringVarP(&aFlag, "flag", "f", "default-val", "Help message for flag")

	// ENV variables
	sentryDSN = os.Getenv("SENTRY_DSN")
}

func setupRemoteLogging() {
	if sentryDSN == "" {
		log.Info("No sentry token specified; error logging not setup")
		return
	}

	sentryTags := map[string]string{
		"buildVersion": BuildVersion,
	}

	hook, err := logrus_sentry.NewWithTagsSentryHook(
		sentryDSN,
		sentryTags,
		[]log.Level{
			// log.DebugLevel,
			// log.InfoLevel,
			log.WarnLevel,
			log.ErrorLevel,
			log.FatalLevel,
			log.PanicLevel,
		},
	)
	if err != nil {
		log.Fatal("Could not setup log monitor")
	}

	hook.Timeout = 5 * time.Second
	hook.StacktraceConfiguration.Enable = true
	log.AddHook(hook)
}

// initConfig only runs if a command runs
func initConfig() {
	if version {
		return // we are only showing version then quitting, no need for anything else
	}
	// Debug mode?
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	setupRemoteLogging()
	log.Debug("Debug mode")
}
