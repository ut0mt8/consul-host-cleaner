package main

import (
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/namsral/flag"

	"consul-route53-sync/internal/consul"
)

type config struct {
	addresses string
	grpc      int
	http      int
	timeout   int
	interval  int
}

func main() {
	var conf config

	log := hclog.New(&hclog.LoggerOptions{
		Name: "cleaner",
	})

	flag.StringVar(&conf.addresses, "consul-addresses", "", "go-netaddrs formated consul servers defintion [REQUIRED]")
	flag.IntVar(&conf.grpc, "consul-grpc-port", 8502, "grpc port of consul server")
	flag.IntVar(&conf.http, "consul-http-port", 8500, "http port of consul server")
	flag.IntVar(&conf.timeout, "consul-http-timeout", 5, "http timeout for connecting to consul server")
	flag.IntVar(&conf.interval, "refresh-interval", 20, "interval between sync")
	flag.Parse()

	if conf.addresses == "" {
		flag.Usage()
		log.Error("required parameters missing")
		return
	}

	cm, err := consul.NewConsulManager(
		conf.addresses,
		consul.WithGRPCPort(conf.grpc),
		consul.WithHTTPPort(conf.http),
		consul.WithTimeout(conf.timeout),
	)
	if err != nil {
		log.Error("create consul manager", "error", hclog.Fmt("%s", err))
		return
	}

	go cm.Run()
	defer cm.Stop()

	for range time.NewTicker(time.Duration(conf.interval) * time.Second).C {
		log.Info("clean fired")

		// first handle left or failed members from gossip
		nodes, err := cm.GetFailedMembers()
		if err != nil {
			log.Error("consul get failed members", "error", hclog.Fmt("%s", err))
			continue
		}

		for _, node := range nodes {
			log.Info("pruning", "member", hclog.Fmt("%s", node))

			err = cm.ForceLeavePrune(node)
			if err != nil {
				log.Error("pruning", "member", hclog.Fmt("%s", node), "error", hclog.Fmt("%s", err))
			}
		}

		// second clean up catalog from reliquas
		nodes, err = cm.GetEmptyNodes()
		if err != nil {
			log.Error("consul get empty nodes", "error", hclog.Fmt("%s", err))
			continue
		}

		for _, node := range nodes {
			log.Info("deregistering from catalog", "node", hclog.Fmt("%s", node))

			err = cm.DeregisterNode(node)
			if err != nil {
				log.Error("deregistering from catalog", "node", hclog.Fmt("%s", node), "error", hclog.Fmt("%s", err))
			}
		}

	}
}
