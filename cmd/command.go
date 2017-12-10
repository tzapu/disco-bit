package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tzapu/disco-bit/core"
)

var aFlag string

// aCmd represents a example command
var aCmd = &cobra.Command{
	Use:   "command",
	Short: "a command",
	Long:  "a example command",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Local flag", aFlag)
		core.StartCommand()
	},
}

func init() {
	RootCmd.AddCommand(aCmd)
	log.Println(os.Getenv("B_KEY"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")
	// NOTICE: there s a few issues with persistant flags
	// https://github.com/spf13/cobra/search?q=PersistentFlags&type=Issues&utf8=✓

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aCmd.Flags().StringVarP(&aFlag, "flag", "f", "default-val", "Help message for flag")
}
