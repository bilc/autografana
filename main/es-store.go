/**********************************************************
 * Author        : biliucheng
 * Email         : bilc_dev@163.com
 * Last modified : 2019-05-22 15:01
 * Filename      : es-store.go
 * Description   : es storage cmd
 * Modified By   :
 * *******************************************************/
package main

import (
	".."
	"encoding/json"
	"flag"
	"fmt"
)

func main() {
	es := flag.String("es", "http://127.0.0.1:9200", "es url")
	doc := flag.String("doc",
		`{"service":"mysql","model":"aa",
	"@timestamp":"2019-05-20T10:00:00Z",
	"FILTER_db":"dbx", "METRIC_qps":100}`,
		"doc")
	flag.Parse()

	cli, err := autografana.NewEsClient(*es)
	if err != nil {
		fmt.Println(err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(*doc), &m)
	if err != nil {
		fmt.Println(err)
	}
	err = cli.PutDoc(m)
	fmt.Println(err)
}
