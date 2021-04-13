package main

import (
	"flag"
	"github.com/Navid2zp/citus-failover/config"
	"github.com/Navid2zp/citus-failover/core"
)

func main() {
	configFile := flag.String("f", "config.yml", "Config file path")
	flag.Parse()
	_, err := config.InitConfig(*configFile)
	if err != nil {
		panic(err)
	}
	core.InitMonitor()
	core.Monitor()
}