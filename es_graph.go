package autografana

import (
	//	"bytes"
	"context"
	//:	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/bilc/grafana-sdk"
	"github.com/olivere/elastic"
)

func FolderUid(service string) string {
	return service
}

func Es2Grafana(esUrl, service, model string, grafanaUrl string, grafanaApiKey string) error {

	index := IndexNameCommon(service, model)
	filters, metrics, err := extractEs(esUrl, index)
	if err != nil {
		return err
	}

	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)

	folderResp, err := grafanaCli.GetFolder(FolderUid(service))
	if err != nil || folderResp.ID == 0 {
		folderResp, err = grafanaCli.CreateFolder(sdk.Folder{UID: FolderUid(service), Title: FolderUid(service)})
	}
	fmt.Println("Folder:", folderResp, err)

	ds := NewEsDataSource(esUrl, index)
	status, err := grafanaCli.CreateDatasource(ds)
	if err != nil {
		return fmt.Errorf("createdatasource err %v", err)
	}
	b, _ := json.Marshal(status)
	fmt.Println("---datasource: ", string(b))
	dashboard := NewGraphBoard(index, filters, metrics, index)
	b, _ = json.Marshal(dashboard)
	fmt.Println("---dashboard: ", string(b))

	resp, err := grafanaCli.SetDashboard(*dashboard, true, folderResp.ID)
	b, _ = json.Marshal(resp)
	fmt.Println("Debug dashboard rsp:", string(b))
	if err != nil {
		return fmt.Errorf("setDashboard err %v", err)
	}
	return nil
}

func extractEs(esUrl, index string) ([]string, []string, error) {
	esCli, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil {
		return nil, nil, err
	}
	exists, err := esCli.IndexExists(index).Do(context.Background())
	if err != nil || !exists {
		return nil, nil, fmt.Errorf("not exists or es client err %v", err)
	}
	reply, err := esCli.IndexGet(index).Do(context.Background())
	if err != nil || len(reply) == 0 {
		return nil, nil, fmt.Errorf("not exists or es client err %v", err)
	}
	var indexInfo *elastic.IndicesGetResponse
	for _, j := range reply {
		indexInfo = j
		break
	}
	filters := make([]string, 0)
	metrics := make([]string, 0)
	if v, ok := indexInfo.Mappings["_doc"]; ok {
		v1 := v.(map[string]interface{})["properties"]
		tmp := v1.(map[string]interface{})
		for field, _ := range tmp {
			if strings.HasPrefix(field, FIELD_FILTER_PREFIX) {
				filters = append(filters, field)
			} else if strings.HasPrefix(field, FIELD_METRIC_PREFIX) {
				metrics = append(metrics, field)
			}
		}
	}
	fmt.Println(filters, metrics)
	return filters, metrics, nil
}

func NewEsDataSource(esUrl string, db string) sdk.Datasource {
	//	jsonData
	//	esVersion: 60
	//keepCookies: []
	//timeField: "@timestamp"
	tmp := false
	ds := sdk.Datasource{
		Access:    "proxy",
		BasicAuth: &tmp,
		Name:      db,
		Database:  &db,
		URL:       esUrl,
		Type:      "elasticsearch",
		JSONData:  map[string]interface{}{"esVersion": "60", "timeField": "@timestamp", "keepCookies": []string{}},
	}
	return ds
}

func NewGraphBoard(myDataSource string, myFilters, myMetrics []string, myTitle string) *sdk.Board {
	var myID uint = 1
	var board sdk.Board
	err := json.Unmarshal([]byte(es_graph_json), &board)
	if err != nil {
		fmt.Println("111", err)
		return nil
	}
	board.Title = myTitle
	board.Annotations.List = board.Annotations.List[0:0]
	templateVar := board.Templating.List[0]
	board.Templating.List = board.Templating.List[0:0]
	templateVar.Datasource = &myDataSource

	luceneQuery := ""
	for _, filter := range myFilters {
		luceneQuery += fmt.Sprintf("%s:$%s AND ", filter, filter)
		templateVar.Label = filter
		templateVar.Name = filter
		templateVar.Query = fmt.Sprintf("{\"find\":\"terms\",\"field\":\"%s\"}", filter)
		//	templateVar.Definition = templateVar.Query
		board.Templating.List = append(board.Templating.List, templateVar)
	}

	if len(luceneQuery) > 5 {
		luceneQuery = luceneQuery[0 : len(luceneQuery)-5]
	}
	//这里使用指针，如果*board.Panels[0]会丢失数据
	panelb, _ := json.Marshal(board.Panels[0])
	//	panelVar := *board.Panels[0]
	board.Panels = board.Panels[0:0]

	for i, metric := range myMetrics {
		var panel sdk.Panel
		//深拷贝
		json.Unmarshal(panelb, &panel)

		panel.Datasource = &myDataSource
		*panel.GridPos.X = (i % 3) * 8
		*panel.GridPos.Y = (i / 3) * 8
		panel.ID = myID
		panel.Title = metric
		myID += 1

		panel.GraphPanel.Targets[0].Metrics[0].Field = metric
		panel.GraphPanel.Targets[0].Query = luceneQuery

		board.Panels = append(board.Panels, &panel)
	}
	return &board
}
