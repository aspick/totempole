package totempole

import (
    "fmt"
    "bytes"
    "strings"
    "os/exec"
    "context"
    "log"
    "github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

var cmdCtx context.Context
var cmdCancel context.CancelFunc

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
    // Start hook
    fmt.Println("service start")
    totemfile, e := readTotemfile()
    if e != nil {
        return e
    }
    cmdCtx, cmdCancel = context.WithCancel(context.Background())
    logger.Info(totemfile)

	go p.run(totemfile)
	return nil
}

func (p *program) run(totemfile Totemfile) {
	// Do work here
    logger.Info("run")
    for _ , v := range totemfile.Daemons {
        go func(){
            for ;; {
                cmdExec := ""
                var cmdArgs []string
                if v.Ps != "" {
                    cmdExec = "powershell"
                    cmdArgs = []string{v.Ps}
                }
                if v.Cmd != "" {
                    cmdExec = "cmd"
                    cmdArgs = []string{v.Cmd}
                }
                if v.Sh != "" {
                    token := strings.Split(v.Sh, " ")
                    cmdExec = token[0]
                    cmdArgs = token[1:]
                }
                cmd := exec.CommandContext(cmdCtx, cmdExec, cmdArgs...)
                cmd.Dir = v.Pwd
                // [todo] pipe setup for log
                var out bytes.Buffer
                cmd.Stdout = &out

                cmdError := cmd.Start()
                if cmdError != nil {
                    logger.Error(cmdError)
                    logger.Info(out.String())
                    break
                }
                logger.Info(v.Name + " is started")
                cmdError = cmd.Wait()
                logger.Info(v.Name + " is stopped")
                if cmdError != nil {
                    logger.Error(cmdError)
                    logger.Info(out.String())
                    break
                }

                select {
                case <- cmdCtx.Done():
                    break
                default:
                    logger.Info(v.Name + " is restarting...")
                }
            }
        }()

    }
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
    // stop hook
    cmdCancel()

	return nil
}

func Service() service.Service {
    svcConfig := &service.Config {
        Name: "Totempole",
        DisplayName: "Totempole",
        Description: "Daemonizer",
        Arguments: []string{"run"},
    }

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

    return s
}
