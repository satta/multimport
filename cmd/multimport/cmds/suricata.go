package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satta/multimport/importer"
	"github.com/satta/multimport/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func suriMain(cmd *cobra.Command, args []string) {
	log.SetFormatter(util.UTCFormatter{
		Formatter: &log.JSONFormatter{},
	})

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
		log.SetFormatter(util.UTCFormatter{
			Formatter: &log.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			},
		})
		log.SetOutput(file)
	}

	logjson, _ := rootCmd.PersistentFlags().GetBool("logjson")
	if logjson {
		log.SetFormatter(util.UTCFormatter{
			Formatter: &log.JSONFormatter{},
		})
	}

	bufSize, _ := rootCmd.PersistentFlags().GetUint("bufsize")
	inChan := make(chan []byte, bufSize)
	var dropped uint64
	var incoming uint64
	var accepted uint64

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.WithFields(log.Fields{
				"domain":          "buffer",
				"buffer-capacity": cap(inChan),
				"buffer-length":   len(inChan),
				"dropped":         dropped,
				"accepted":        accepted,
				"received":        incoming,
			}).Info()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	discard, _ := rootCmd.PersistentFlags().GetBool("discard")
	if discard {
		go func(in chan []byte) {
			for v := range in {
				_ = v
			}

		}(inChan)
	} else {
		nofJobs, _ := rootCmd.PersistentFlags().GetUint("jobs")
		vastPath, _ := rootCmd.PersistentFlags().GetString("vast-path")
		vastParams, _ := rootCmd.PersistentFlags().GetStringSlice("extra-params")
		vastBufSize, _ := rootCmd.PersistentFlags().GetInt("vastbufsize")
		log.Debugf("starting %d jobs", nofJobs)
		for i := uint(0); i < nofJobs; i++ {
			importer := importer.MakeImporter(inChan, fmt.Sprintf("suri_%d", i), vastPath, vastParams)
			importer.SetBufSize(vastBufSize)
			go importer.Run("suricata", ctx)
		}
	}	

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Info("signal mask set up")
	go func() {
		for sig := range c {
			if sig == syscall.SIGTERM || sig == syscall.SIGINT || sig == syscall.SIGKILL {
				log.Info("received signal %v, terminating", sig)
				cancel()
				os.Exit(1)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		incoming++
		select {
		case inChan <- []byte(line):
			accepted++
		default:
			dropped++
		}
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
