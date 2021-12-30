package main

import (
	"dyip-sync/client"
	"dyip-sync/config"
	"dyip-sync/meta"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	configFileClient     string
	generateConfigClient bool
	printUsageClient     bool
)

func init() {
	flag.StringVar(&configFileClient, "c", "client.conf", "config file path, default client.conf")
	flag.BoolVar(&generateConfigClient, "g", false, "generate config, default client.conf")
	flag.BoolVar(&printUsageClient, "h", false, "print usage")

	flag.Usage = clientUsage
}

func main() {
	flag.Parse()

	if printUsageClient {
		clientUsage()
		os.Exit(0)
	}

	if generateConfigClient {
		metaData := meta.ClientMeta{}
		metaData.Generate()
		err := config.WriteConfig(configFileClient, &metaData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	metaData := meta.ClientMeta{}
	err := config.ReadConfig(configFileClient, &metaData)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if metaData.Empty() {
		_, _ = fmt.Fprintln(os.Stderr, "config file error: empty")
		os.Exit(1)
	}

	client.MetaData = metaData

	ticker := time.NewTicker(time.Duration(metaData.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err = client.Sync()

		timeStr := time.Now().Format("2006-01-02 15:04:05")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s sync error: %v\n", timeStr, err)
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "%s sync success\n", timeStr)
		}
	}

}

func clientUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:")
	_, _ = fmt.Fprintln(os.Stderr, "  client startup: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-client [-c config file]")
	_, _ = fmt.Fprintln(os.Stderr, "  client startup in background: ")
	_, _ = fmt.Fprintln(os.Stderr, "    nohup dyip-client [-c config file] &")
	_, _ = fmt.Fprintln(os.Stderr, "  generate demo config file: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-client -g [-c config file]")
	_, _ = fmt.Fprintln(os.Stderr, "  print usage: ")
	_, _ = fmt.Fprintln(os.Stderr, "    dyip-client -h")
	_, _ = fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
}
