package importer

import (
	"bufio"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Importer struct {
	name string
}

func (i *Importer) Run(inChan chan []byte, importType string, name string) error {
	i.name = name
	for {
		stopChan := make(chan bool)
		cmd := exec.Command("vast", "import", importType)
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
			for line := range inChan {
				select {
				case <-stopChan:
					return
				default:
					stdin.Write(line)
				}

			}
		}()
		go func() {
			log.Debug(i.name + " started stderr scanner goroutine")
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(i.name + ": " + scanner.Text())
			}
			log.Debug(i.name + " closed stderr scanner goroutine")
		}()
		err = cmd.Wait()
		if err != nil {
			log.Errorf("importer finished with error: %v", err)
		}
		close(stopChan)
	}
}
