package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/satta/multimport/importer"
	"github.com/satta/multimport/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func suriMain(cmd *cobra.Command, args []string) {
	log.SetFormatter(util.UTCFormatter{&log.JSONFormatter{}})

	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	if verbose {
		log.Debug("verbose log output enabled")
		log.SetLevel(log.DebugLevel)
	}

	logfilename, _ := rootCmd.PersistentFlags().GetString("logfile")
	if len(logfilename) > 0 {
		log.Println("Switching to log file", logfilename)
		file, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetFormatter(util.UTCFormatter{&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		}})
		log.SetOutput(file)
	}

	logjson, _ := rootCmd.PersistentFlags().GetBool("logjson")
	if logjson {
		log.SetFormatter(util.UTCFormatter{&log.JSONFormatter{}})
	}

	bufSize, _ := rootCmd.PersistentFlags().GetUint("bufsize")
	inChan := make(chan []byte, bufSize)
	dropped := 0

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.WithFields(log.Fields{
				"domain":          "buffer",
				"buffer-capacity": cap(inChan),
				"buffer-length":   len(inChan),
				"dropped":         dropped,
			}).Info()
		}
	}()

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
		select {
		case inChan <- []byte(line):
		default:
			log.Debug("channel full, discarding line")
			dropped++
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
