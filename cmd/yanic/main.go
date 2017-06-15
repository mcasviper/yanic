package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mcasviper/yanic/database"
	"github.com/mcasviper/yanic/database/all"
	"github.com/mcasviper/yanic/meshviewer"
	"github.com/mcasviper/yanic/respond"
	"github.com/mcasviper/yanic/rrd"
	"github.com/mcasviper/yanic/runtime"
	"github.com/mcasviper/yanic/webserver"
)

var (
	configFile  string
	config      *runtime.Config
	collector   *respond.Collector
	connections database.Connection
	nodes       *runtime.Nodes
)

func main() {
	var importPath string
	var timestamps bool
	flag.StringVar(&importPath, "import", "", "import global statistics from the given RRD file, requires influxdb")
	flag.StringVar(&configFile, "config", "config.toml", "path of configuration file (default:config.yaml)")
	flag.BoolVar(&timestamps, "timestamps", true, "print timestamps in output")
	flag.Parse()

	if !timestamps {
		log.SetFlags(0)
	}
	log.Println("Yanic say hello")

	config, err := runtime.ReadConfigFile(configFile)
	if err != nil {
		panic(err)
	}

	connections, err = all.Connect(config.Database.Connection)
	if err != nil {
		panic(err)
	}
	database.Start(connections, config)
	defer database.Close(connections)

	if connections != nil && importPath != "" {
		importRRD(importPath)
		return
	}

	nodes = runtime.NewNodes(config)
	nodes.Start()
	meshviewer.Start(config, nodes)

	if config.Webserver.Enable {
		log.Println("starting webserver on", config.Webserver.Bind)
		srv := webserver.New(config.Webserver.Bind, config.Webserver.Webroot)
		go srv.Close()
	}

	if config.Respondd.Enable {
		// Delaying startup to start at a multiple of `duration` since the zero time.
		if duration := config.Respondd.Synchronize.Duration; duration > 0 {
			now := time.Now()
			delay := duration - now.Sub(now.Truncate(duration))
			log.Printf("delaying %0.1f seconds", delay.Seconds())
			time.Sleep(delay)
		}

		collector = respond.NewCollector(connections, nodes, config.Respondd.Interface, config.Respondd.Port)
		collector.Start(config.Respondd.CollectInterval.Duration)
		defer collector.Close()
	}

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)
}

func importRRD(path string) {
	log.Println("importing RRD from", path)
	for ds := range rrd.Read(path) {
		connections.InsertGlobals(
			&runtime.GlobalStats{
				Nodes:   uint32(ds.Nodes),
				Clients: uint32(ds.Clients),
			},
			ds.Time,
		)
	}
}
