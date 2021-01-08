package cmd

// DCSO FEVER
// Copyright (c) 2018, DCSO GmbH

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
