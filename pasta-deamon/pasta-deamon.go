package main

import (
	"github.com/takama/daemon"
	"github.com/egomezbpedro/cli-pasta/clipboard"
    "github.com/egomezbpedro/cli-pasta/database"
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "sync"
)

const (
    // name of the service
    name        = "CliPasta"
    description = "CliPasta is a service to store and retrieve text snippets from the clipboard"
)
var (
    stdlog, errlog *log.Logger
    
    d = database.Database{
        BucketName: "data",
        DatabaseName: "/usr/local/var/pasta.db",
    }
    clip = clipboard.Clipboard{}
    wg = sync.WaitGroup{}
)

// Service has embedded daemon
type Service struct {
    daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

    usage := "Usage: CliPasta install | remove | start | stop | status"

    // if received any kind of command, do it
    if len(os.Args) > 1 {
        command := os.Args[1]
        switch command {
        case "install":
            return service.Install()
        case "remove":
            return service.Remove()
        case "start":
            return service.Start()
        case "stop":
            return service.Stop()
        case "status":
            return service.Status()
        default:
            return usage, nil
        }
    }

    wg.Add(1)
    go func() {
		for {
            stdlog.Println("Writing to bucket")
			d.WriteToBucket(d.DatabaseName, d.BucketName, clip.WatchClipboard())
		}
	}()
    wg.Wait()
	
	interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    // loop work cycle with accept connections or interrupt
    // by system signal
    for {
        select {
        case killSignal := <-interrupt:
            stdlog.Println("Got signal:", killSignal)
            if killSignal == os.Interrupt {
                return "Daemon was interruped by system signal", nil
            }
            return "Daemon was killed", nil
		default:
			return "usage", nil
        }	
    }
}

func init() {
    stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {

    d.CreateBucket(d.DatabaseName, d.BucketName)
    
    srv, err := daemon.New(name, description, daemon.UserAgent, nil...)
    if err != nil {
        errlog.Println("Error: ", err)
        os.Exit(1)
    }
    service := &Service{srv}
    status, err := service.Manage()
    if err != nil {
        errlog.Println(status, "\nError: ", err)
        os.Exit(1)
    }
    fmt.Println(status)
}
