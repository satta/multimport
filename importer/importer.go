package importer

import (
	"bufio"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Importer struct {
	InChan   chan []byte
	VastPath string
	Logger   *log.Entry
	Params   []string
}

func MakeImporter(inChan chan []byte, name string, vastPath string, params []string) *Importer {
	i := &Importer{
		InChan: inChan,
		Logger: log.WithFields(log.Fields{
			"importer": name,
		}),
		VastPath: vastPath,
		Params:   params,
	}
	return i
}

func (i *Importer) Run(importType string) error {
	for {
		stopChan := make(chan bool)
		params := append(i.Params, "import", importType)
		i.Logger.Debugf("starting command '%s' with params %v", i.VastPath, params)
		cmd := exec.Command(i.VastPath, params...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for line := range i.InChan {
				select {
				case <-stopChan:
					return
				default:
					stdin.Write(line)
				}

			}
		}()
		go func() {
			i.Logger.Debug("started stderr scanner goroutine")
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				i.Logger.Info(scanner.Text())
			}
			i.Logger.Debug("closed stderr scanner goroutine")
		}()
		err = cmd.Wait()
		if err != nil {
			i.Logger.Errorf("importer finished with error: %v", err)
		}
		close(stopChan)
	}
}
