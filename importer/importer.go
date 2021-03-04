package importer

import (
	"bufio"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Importer struct {
	InChan            chan []byte
	VastPath          string
	Logger            *log.Entry
	Params            []string
	Count             uint64
	CountLock         sync.Mutex
	CountIntervalSecs uint
}

func MakeImporter(inChan chan []byte, name string, vastPath string, params []string) *Importer {
	i := &Importer{
		InChan: inChan,
		Logger: log.WithFields(log.Fields{
			"domain":   "importer",
			"importer": name,
		}),
		VastPath:          vastPath,
		Params:            params,
		CountIntervalSecs: 10,
	}
	return i
}

func (i *Importer) Run(importType string) error {
	i.Logger.Debugf("importer started")
	for {
		stopChan := make(chan bool, 2)
		params := append(i.Params, "import", importType)
		i.Logger.Debugf("starting command '%s' with params %v", i.VastPath, params)
		cmd := exec.Command(i.VastPath, params...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		stdinWriter := bufio.NewWriter(stdin)
		go func(sChan chan bool) {
			i.Logger.Debug("started stdin writer goroutine")
			defer i.Logger.Debug("left stdin writer goroutine")
			for {
				for line := range i.InChan {
					select {
					case <-sChan:
						return
					default:
						i.CountLock.Lock()
						i.Count++
						i.CountLock.Unlock()
						_, err := stdinWriter.Write(line)
						if err != nil {
							i.Logger.Debugf("could not write: %s", error(err))
						}
					}
				}
			}
		}(stopChan)
		stdoutReader := bufio.NewReader(stdout)
		go func(sChan chan bool) {
			i.Logger.Debug("started stdout scanner goroutine")
			defer i.Logger.Debug("left stdout scanner goroutine")
			for {
				select {
				case <-sChan:
					return
				default:
					i.Logger.Debug("started stdout scanner loop")
					scanner := bufio.NewScanner(stdoutReader)
					for scanner.Scan() {
						select {
						case <-sChan:
							return
						default:
							i.Logger.Info(scanner.Text())
						}
					}
					if err := scanner.Err(); err != nil {
						i.Logger.Errorf("error reading stdout: %s", err)
					}
					i.Logger.Debug("end of stdout scanner loop")
				}
			}
		}(stopChan)
		stderrReader := bufio.NewReader(stderr)
		go func(sChan chan bool) {
			i.Logger.Debug("started stderr scanner goroutine")
			defer i.Logger.Debug("left stderr scanner goroutine")
			for {
				select {
				case <-sChan:
					return
				default:
					i.Logger.Debug("started stderr scanner loop")
					scanner := bufio.NewScanner(stderrReader)
					for scanner.Scan() {
						select {
						case <-sChan:
							return
						default:
							i.Logger.Info(scanner.Text())
						}
					}
					if err := scanner.Err(); err != nil {
						i.Logger.Errorf("error reading stderr: %s", err)
					}
					i.Logger.Debug("end of stderr scanner loop")
				}
			}
		}(stopChan)
		go func(sChan chan bool) {
			i.Logger.Debug("started logger goroutine")
			defer i.Logger.Debug("left logger goroutine")
			for {
				select {
				case <-sChan:
					return
				case <-time.After(time.Duration(i.CountIntervalSecs) * time.Second):
					i.CountLock.Lock()
					myCount := i.Count
					i.Count = 0
					i.CountLock.Unlock()
					i.Logger.Infof("processed %f lines per second (%d total)", float64(myCount)/float64(i.CountIntervalSecs), myCount)
				}
			}
		}(stopChan)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			i.Logger.Errorf("importer finished with error: %v", err)
		}
		close(stopChan)
		time.Sleep(10 * time.Second)
	}
}
