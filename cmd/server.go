package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/thecasualcoder/dobby/pkg/handler"
	"github.com/thecasualcoder/dobby/pkg/server"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func serverCmd() cli.Command {
	return cli.Command{
		Name:        "server",
		Description: "run dobby in server mode",
		Usage:       "run dobby server",
		Action:      runServer,
		Flags:       serverFlags(),
	}
}

func serverFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "bind-address, a",
			Value:  "127.0.0.1",
			Usage:  "Address of dobby server.",
			EnvVar: "BIND_ADDR",
		},
		cli.StringFlag{
			Name:   "port, p",
			Value:  "4444",
			Usage:  "Port which will be used by dobby server.",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "initial-health",
			EnvVar: "INITIAL_HEALTH",
			Usage:  "Sets the Initial health of the server (/health) (true|false)",
			Value:  "true",
		},
		cli.StringFlag{
			Name:   "initial-readiness",
			EnvVar: "INITIAL_READINESS",
			Usage:  "Sets the Initial readiness of the server (/readiness) (true|false)",
			Value:  "true",
		},
		cli.Int64Flag{
			Name:   "initial-delay",
			EnvVar: "INITIAL_DELAY",
			Usage:  "Sets the Initial delay to start the server (in seconds)",
			Value:  0,
		},
		cli.StringFlag{
			Name:   "proxy-config-path",
			EnvVar: "PROXY_CONFIG_PATH",
			Usage:  "Sets the proxy configuration path",
			Value:  "",
		},
	}
}

func watchForProxyChanges(watcher *fsnotify.Watcher, pathToWatch string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var fileExists = func(file string) bool {
		fileFullPath := filepath.Join(wd, file)
		if _, err := os.Stat(fileFullPath); err != nil {
			return false
		}
		return true
	}

	var addProxies = func(file string) error {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		proxyRequests, err := handler.NewProxyRequests(bytes)
		if err != nil {
			return err
		}

		handler.AddProxies(file, proxyRequests)
		return nil
	}

	var removeProxies = func(file string) {
		fileFullPath := filepath.Join(wd, file)
		handler.RemoveProxies(fileFullPath)
	}

	var addExistingProxy = func() {
		files, err := filepath.Glob(filepath.Join(wd, pathToWatch, "*.yaml"))
		if err != nil {
			return
		}
		for _, file := range files {
			err := addProxies(file)
			if err != nil {
				log.Println(fmt.Sprintf("unable to add proxies from '%s', dynamically", file), err)
			}
		}
	}

	go func() {
		addExistingProxy()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					if fileExists(event.Name) {
						fileFullPath := filepath.Join(wd, event.Name)
						err := addProxies(fileFullPath)
						if err != nil {
							log.Println("unable to add proxies dynamically", err)
						} else {
							log.Println("configured proxies from: ", event.Name)
						}
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					removeProxies(event.Name)
					log.Println("removed proxies from", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	watchPath := filepath.Join(wd, pathToWatch)
	log.Println(fmt.Sprintf("started to watch %s for proxy changes", watchPath))

	return watcher.Add(pathToWatch)
}

func pathToWatch(context *cli.Context) (string, bool) {
	proxyPath := context.String("proxy-config-path")
	if _, err := os.Stat(proxyPath); err != nil {
		return "", false
	}
	return proxyPath, true
}

func runServer(context *cli.Context) {
	bindAddress := context.String("bind-address")
	port := context.String("port")

	initialDelay := time.Duration(context.Int64("initial-delay")) * time.Second
	time.Sleep(initialDelay)

	initialHealth := true
	if health, err := strconv.ParseBool(context.String("initial-health")); err == nil {
		initialHealth = health
	}

	initialReadiness := true
	if readiness, err := strconv.ParseBool(context.String("initial-readiness")); err == nil {
		initialReadiness = readiness
	}

	if watchDir, ok := pathToWatch(context); ok {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			dieIf(err)
		}

		defer func() {
			err := watcher.Close()
			if err != nil {
				log.Printf("unable to close watcher, reason %v", err)
			}
		}()

		err = watchForProxyChanges(watcher, watchDir)
		if err != nil {
			dieIf(err)
		}
	}

	err := server.Run(bindAddress, port, initialHealth, initialReadiness)
	dieIf(err)
}
