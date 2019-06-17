/**********************************************************
 * Author        : biliucheng
 * Email         : bilc_dev@163.com
 * Last modified : 2019-05-22 15:01
 * Filename      : es2grafana.go
 * Description   : generate grafana dashboard cmd
 * Modified By   :
 * *******************************************************/
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
	//panel := flag.String("panel","graph","matric panel type: graph, heatmap or all")
	flag.Parse()
	//"eyJrIjoidmphWWxFSVg1UzdXMXV3T1hoNWcwVFd2alp6NUQxd2siLCJuIjoiYXBpa2V5Y3VybCIsImlkIjoxfQ=="
	panel := make(map[string][]string)
	panel["METRIC_bill"] = []string{"graph"}
	panel["METRIC_qps"] = []string{"heatmap"}
	tagsSorts := []string{"TAG_region", "TAG_az", "TAG_host", "TAG_source_type", "TAG_flavor", "TAG_user"}

	err := autografana.Es2Grafana(*es, *service, *model, *grafana, *key, nil, tagsSorts, panel)
	fmt.Println(err)

	temp := autografana.TemplateVars{
		{Name: "TAG_user", Sort: 4},
		{Name: "TAG_b", Sort: 5},
		{Name: "TAG_host", Sort: 3},
		{Name: "TAG_AZ", Sort: 2},
		{Name: "TAG_region", Sort: 1},
		{Name: "TAG_a", Sort: 6},
		{Name: "TAG_c", Sort: 7},
	}
	autografana.SortTemplatingList(temp)
	for _, t := range temp {
		fmt.Println(t)
	}
}
