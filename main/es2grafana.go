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
	// alpha-test key:
	key := "eyJrIjoiSk5ybDV5Tlp6YkZoQU0wUEZmaktHVE52MllVN1oyR0wiLCJuIjoiYWRtaW4iLCJpZCI6MX0="
	// hb api key
	//key := "eyJrIjoiYjV2RUd3cUplTktiZGo4ZGlHYUpFU2VPWGRCWEFJY0UiLCJuIjoiYWRtaW4iLCJpZCI6MX0="
	// stag key:
	//key := "eyJrIjoiNnBoZHJ3ODFGM3pxMjZmeU9qdGZjN05KYzRnT0Z3MXUiLCJuIjoiQWRtaW4iLCJpZCI6MX0="

	flag.Parse()
	mypanel := []autografana.MyPanel{
		{
			Title: "bill used",
			MyMetrics: []autografana.MyMetric{
				{Field: "SUM_METRIC_bill_current", Type: "sum"},
				{Field: "SUM_METRIC_bill_total", Type: "sum"},
			},
			Type: autografana.PANEL_GRAPH,
			Interval: "10s",
		},{
			Title: "qps used",
			MyMetrics: []autografana.MyMetric{
				{Field: "METRIC_qps_current", Type: "avg"},
				{Field: "METRIC_qps_total", Type: "avg"},
			},
			Type: autografana.PAENL_HEATMAP,
			Interval: "20s",
		},
	}

	tagsSorts := []string{"TAG_region", "TAG_az", "TAG_flavor","TAG_kind","TAG_resource_name","TAG_release","TAG_pin", "TAG_instance_id"}
	tagsCascade := make(map[string][]string)
	tagsCascade["TAG_az"] = []string{"TAG_region"}
	tagsCascade["TAG_resource_name"] = []string{"TAG_instance_id"}
	tagsCascade["TAG_instance_id"] = []string{"TAG_pin","TAG_az","TAG_flavor", "TAG_release"}

	url, err := autografana.Es2Grafana(*es, *esHasAuth, *user, *password, *service, *model, *grafana, key, nil, tagsSorts, tagsCascade, mypanel)
	fmt.Println(url, err)

	allfolder, err := autografana.ListFolders(*grafana, key)
	fmt.Printf("all folder: %+v",allfolder)
	fmt.Println(err)

	uid, err := autografana.GetFolderUid(*grafana, key, "es")
	fmt.Printf("es uid is %+s",uid)

	esfolder, err := autografana.GetAllDashboardInFolder(*grafana, key, "es")
	fmt.Printf("es folder: %+v", esfolder)

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
