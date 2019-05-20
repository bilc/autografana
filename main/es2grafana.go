package main

import (
	".."
	"flag"
	"fmt"
)

func main() {
	es := flag.String("es", "http://127.0.0.1:9200", "es url")
	//	index := flag.String("index", "test1", "es index")
	service := flag.String("service", "mysql", "es index field service")
	model := flag.String("model", "qps", "es index field model")

	grafana := flag.String("grafana", "http://127.0.0.1:3000", "grafana url")
	key := flag.String("key", "eyJrIjoidmphWWxFSVg1UzdXMXV3T1hoNWcwVFd2alp6NUQxd2siLCJuIjoiYXBpa2V5Y3VybCIsImlkIjoxfQ==", "grafana api key")
	flag.Parse()
	//"eyJrIjoidmphWWxFSVg1UzdXMXV3T1hoNWcwVFd2alp6NUQxd2siLCJuIjoiYXBpa2V5Y3VybCIsImlkIjoxfQ=="
	err := autografana.Es2Grafana(*es, *service, *model, *grafana, *key)
	fmt.Println(err)
}
