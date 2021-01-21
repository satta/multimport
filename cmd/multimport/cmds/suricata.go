package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/satta/multimport/importer"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func suriMain(cmd *cobra.Command, args []string) {
	log.SetFormatter(util.UTCFormatter{&log.JSONFormatter{}})
	inChan := make(chan []byte)

	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	if verbose {
		log.Debug("verbose log output enabled")
		log.SetLevel(log.DebugLevel)
	}

	nofJobs, _ := rootCmd.PersistentFlags().GetUint("jobs")
	vastPath, _ := rootCmd.PersistentFlags().GetString("vast-path")
	vastParams, _ := rootCmd.PersistentFlags().GetStringSlice("extra-params")
	log.Debugf("starting %d jobs", nofJobs)
	for i := uint(0); i < nofJobs; i++ {
		importer := importer.MakeImporter(inChan, fmt.Sprintf("suri_%d", i), vastPath, vastParams)
		go importer.Run("suricata")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		inChan <- []byte(line)
	}
}

var runCmd = &cobra.Command{
	Use:   "suricata",
	Short: "import Suricata data",
	Run:   suriMain,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
