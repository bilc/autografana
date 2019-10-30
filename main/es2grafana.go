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
	es := flag.String("es", "http://10.226.134.46:9200", "es url")
	esHasAuth := flag.String("esauth", "http://10.226.134.46:9200", "es has auth url")
	user := flag.String("user", "jvessel", "es user")
	password := flag.String("password", "jvessel-es", "es password")
	service := flag.String("service", "smoke", "es index field service")
	model := flag.String("model", "es", "es index field model")
	grafana := flag.String("grafana", "http://10.226.134.46:3000", "grafana url")
	key := flag.String("key", "admin", "grafana api key")
	flag.Parse()

	/*mypanel := []autografana.MyPanel{
		{
			Title: "bill used",
			Metrics: []string{"SUM_METRIC_bill_current", "SUM_METRIC_bill_total"},
			Type: autografana.PANEL_GRAPH,
			Interval: "10s",
		},{
			Title: "qps used",
			Metrics: []string{"METRIC_qps_current", "METRIC_qps_total"},
			Type: autografana.PAENL_HEATMAP,
			Interval: "20s",
		},
	}*/

	url, err := autografana.Es2Grafana(*es, *esHasAuth, *user, *password, *service, *model, *grafana, key, nil, nil, nil, nil)
	fmt.Println(url, err)
}
