package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multimport",
	Short: "VAST multiprocessed importer",
	Long:  `Helps with importing lots of events in parallel.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().UintP("jobs", "j", 4, "amount of parallel VAST import processes")
	rootCmd.PersistentFlags().UintP("bufsize", "b", 100000, "input event buffer size in EVE lines")
	rootCmd.PersistentFlags().IntP("vastbufsize", "s", 1*1024*1024, "VAST process communication buffer size in bytes")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")
	rootCmd.PersistentFlags().BoolP("discard", "", false, "discard mode, only reads from source without importing (for debugging)")
	rootCmd.PersistentFlags().StringP("vast-path", "", "vast", "VAST executable")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "logfile (stderr if empty)")
	rootCmd.PersistentFlags().BoolP("logjson", "", false, "log in JSON format")
	rootCmd.PersistentFlags().StringSliceP("extra-params", "p", []string{}, "extra parameters to pass to 'vast import', separated by commas")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
