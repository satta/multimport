package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

func versionMain(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version number",
	Run:   versionMain,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
