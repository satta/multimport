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
	inChan := make(chan []byte)

	nofJobs, _ := rootCmd.PersistentFlags().GetUint("jobs")
	vastPath, _ := rootCmd.PersistentFlags().GetString("vastpath")
	log.Debugf("starting %d jobs", nofJobs)
	for i := uint(0); i < nofJobs; i++ {
		importer := importer.MakeImporter(inChan, fmt.Sprintf("suri_%d", i), vastPath)
		go importer.Run("suricata")
		log.Debugf("importer %d started", i)
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
