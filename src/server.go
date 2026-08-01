package main

import (
	"dyip-sync/src/config"
	"dyip-sync/src/frontend"
	"dyip-sync/src/meta"
	"dyip-sync/src/server"
	"flag"
	"fmt"
	"os"

	"github.com/savsgio/atreugo/v11"
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

	atreugoServer.StaticCustom("/", &atreugo.StaticFS{
		FS:              frontend.FS,
		IndexNames:      []string{"index.html"},
		AllowEmptyRoot:  true,
		Compress:        true,
		CompressBrotli:  true,
		AcceptByteRange: true,
		// Browsers always send Accept-Encoding, which sends fasthttp down its
		// on-demand-compression path. For the bare "/" that path resolves to an
		// empty file path and embed.FS.Open("") fails with "invalid argument",
		// so GET / returns 404 from a browser (curl hides it by omitting
		// Accept-Encoding). Rewrite "/" to "/index.html" so the handler looks
		// up a real file that compresses/serves normally. All other paths are
		// left untouched.
		PathRewrite: func(ctx *atreugo.RequestCtx) []byte {
			if string(ctx.Path()) == "/" {
				return []byte("/index.html")
			}
			return ctx.Path()
		},
	})

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
