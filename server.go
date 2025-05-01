package main

import (
	"dyip-sync/config"
	"dyip-sync/meta"
	"dyip-sync/server"
	"flag"
	"fmt"
	"github.com/savsgio/atreugo/v11"
	"os"
)

var (
	configFileServer     string
	generateConfigServer bool
	printUsageServer     bool
)

func init() {
	flag.StringVar(&configFileServer, "c", "server.conf", "config file path, default server.conf")
	flag.BoolVar(&generateConfigServer, "g", false, "generate config, default server.conf")
	flag.BoolVar(&printUsageServer, "h", false, "print usage")

	flag.Usage = serverUsage
}

func main() {
	flag.Parse()

	if printUsageServer {
		serverUsage()
		os.Exit(0)
	}

	if generateConfigServer {
		metaData := meta.ServerMeta{}
		metaData.Generate()
		err := config.WriteConfig(configFileServer, &metaData)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	metaData := meta.ServerMeta{}
	err := config.ReadConfig(configFileServer, &metaData)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if metaData.Empty() {
		_, _ = fmt.Fprintln(os.Stderr, "config file error: empty")
		os.Exit(1)
	}

	metaData.GenerateIpm()

	server.MetaData = metaData
	server.ConfigFileServer = configFileServer

	config := atreugo.Config{
		Addr: fmt.Sprintf("%s:%d", metaData.Address, metaData.Port),
	}
	atreugoServer := atreugo.New(config)

	atreugoServer.GET("/", server.IndexHandler)
	atreugoServer.GET("/sync", server.SyncHandler)
	atreugoServer.GET("/load", server.LoadHandler)
	atreugoServer.GET("/ip", server.IpHandler)

	err = atreugoServer.ListenAndServe()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "start server error: %v\n", err)
		os.Exit(1)
	}
}

func serverUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:")
	_, _ = fmt.Fprintln(os.Stderr, "  server startup: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-server [-c config file]")
	_, _ = fmt.Fprintln(os.Stderr, "  server startup in background: ")
	_, _ = fmt.Fprintln(os.Stderr, "    nohup dyip-server [-c config file] &")
	_, _ = fmt.Fprintln(os.Stderr, "  generate demo config file: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-server -g [-c config file]")
	_, _ = fmt.Fprintln(os.Stderr, "  print usage: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-server -h")
	_, _ = fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
}
