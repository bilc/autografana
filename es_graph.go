/**********************************************************
 * Author        : biliucheng
 * Email         : bilc_dev@163.com
 * Last modified : 2019-05-22 12:59
 * Filename      : es_graph.go
 * Description   : make dashboard by es index
 * Modified By   :
 * *******************************************************/
package autografana

import (
	//	"bytes"
	"context"
	//:	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/bilc/grafana-sdk"
	"gopkg.in/olivere/elastic.v6"
)

const PANEL_GRAPH = "graph"
const PAENL_HEATMAP = "heatmap"

func Es2Grafana(esUrl, service, model string, grafanaUrl string, grafanaApiKey string, gratags []string, panel map[string][]string) error {

	index := IndexNameCommon(service, model)
	tags, metrics, err := ExtractEs(esUrl, index)
	if err != nil {
		return err
	}
	if gratags != nil {
		for _, j := range gratags {
			exist := false
			for _, k := range tags {
				if k == j {
					exist = true
					break
				}
			}
			if !exist {
				return fmt.Errorf("tag %v not exist", j)
			}
		}
		tags = gratags
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
	dashboard := NewGraphBoard(index, tags, metrics, panel, model)
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

func ListServiceModelByExtractEs(esUrl, index string) ([]string,error) {
	esCli, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil {
		return nil, err
	}
	exists, err := esCli.IndexExists(index).Do(context.Background())
	if err != nil || !exists {
		return nil, fmt.Errorf("not exists or es client err %v", err)
	}
	reply, err := esCli.IndexGet(index).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("not exists or es client err %v", err)
	}
	models := make([]string,0)
	prefix := strings.TrimRight(index,"*")+"-"
	for replyKey, _ := range reply {
		if strings.HasPrefix(replyKey, prefix) {
			modelTail := strings.TrimPrefix(replyKey, prefix)
			lastIndex := strings.LastIndex(modelTail, "-")
			model := modelTail[0: lastIndex]
			models = append(models,model)
		}
	}

	fmt.Println("model ",RemoveRepeatedElement(models))
	return RemoveRepeatedElement(models), nil
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func ExtractEs(esUrl, index string) ([]string, []string, error) {
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
	tags := make([]string, 0)
	metrics := make([]string, 0)
	if v, ok := indexInfo.Mappings["_doc"]; ok {
		v1 := v.(map[string]interface{})["properties"]
		tmp := v1.(map[string]interface{})
		for field, _ := range tmp {
			if strings.HasPrefix(field, FIELD_TAG_PREFIX) {
				tags = append(tags, field)
			} else if strings.HasPrefix(field, FIELD_METRIC_PREFIX) {
				metrics = append(metrics, field)
			}
		}
	}
	sort.Strings(tags)
	sort.Strings(metrics)
	fmt.Println(tags, metrics)
	return tags, metrics, nil
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

func NewGraphBoard(myDataSource string, mytags, myMetrics []string, panel map[string][]string, myTitle string) *sdk.Board {
	var myID uint = 1
	var board sdk.Board
	err := json.Unmarshal([]byte(es_grafana_json), &board)
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
	for _, tag := range mytags {
		luceneQuery += fmt.Sprintf("%s:$%s AND ", tag, tag)
		templateVar.Label = tag
		templateVar.Name = tag
		templateVar.Query = fmt.Sprintf("{\"find\":\"terms\",\"field\":\"%s\"}", tag)
		//	templateVar.Definition = templateVar.Query
		board.Templating.List = append(board.Templating.List, templateVar)
	}

	if len(luceneQuery) > 5 {
		luceneQuery = luceneQuery[0 : len(luceneQuery)-5]
	}
	//这里使用指针，如果*board.Panels[0]会丢失数据
	//panelb, _ := json.Marshal(board.Panels[0])
	//	panelVar := *board.Panels[0]
	board.Panels = board.Panels[0:0]
	for i, metric := range myMetrics {
		var panelMatic *sdk.Panel
		if _, ok := panel[metric]; ok {
			panelTypes := panel[metric]
			for _, panelType := range panelTypes {
				if panelType == PANEL_GRAPH{
					panelMatic = NewGraphPanel(myDataSource, myID, i, metric, luceneQuery)
				}else if panelType == PAENL_HEATMAP{
					panelMatic = NewHeatmapPanel(myDataSource, myID, i, metric, luceneQuery)
				}
				myID += 1
				board.Panels = append(board.Panels, panelMatic)
			}
		}else{
			panelMatic = NewGraphPanel(myDataSource, myID, i, metric, luceneQuery)
			myID += 1
			board.Panels = append(board.Panels, panelMatic)
		}
	}
	return &board
}

func FolderUid(service string) string {
	return service
}

func NewGraphPanel(myDataSource string, panelId uint, metrixIndex int, metric string, luceneQuery string) *sdk.Panel {
	var graphPanel sdk.Panel
	err := json.Unmarshal([]byte(graph_panel_json), &graphPanel)
	if err != nil {
		fmt.Println("unmarshl graph panel json error:", err)
		return nil
	}

	graphPanel.Datasource = &myDataSource
	*graphPanel.GridPos.X = (metrixIndex % 3) * 8
	*graphPanel.GridPos.Y = (metrixIndex / 3) * 8
	graphPanel.ID = panelId
	graphPanel.Title = metric

	graphPanel.GraphPanel.Targets[0].Metrics[0].Field = metric
	graphPanel.GraphPanel.Targets[0].Query = luceneQuery

	return &graphPanel
}

func NewHeatmapPanel(myDataSource string, panelId uint, metrixIndex int, metric string, luceneQuery string) *sdk.Panel {
	var heatmapPanel sdk.Panel
	err := json.Unmarshal([]byte(heatmap_panel_json), &heatmapPanel)
	if err != nil {
		fmt.Println("unmarshl heatmap panel json error:", err)
		return nil
	}

	heatmapPanel.Datasource = &myDataSource
	*heatmapPanel.GridPos.X = (metrixIndex % 3) * 8
	*heatmapPanel.GridPos.Y = (metrixIndex / 3) * 8
	heatmapPanel.ID = panelId
	heatmapPanel.Title = metric
	heatmapPanel.HeatmapPanel.Targets[0].Metrics[0].Field = metric
	heatmapPanel.HeatmapPanel.Targets[0].Query = luceneQuery

	return &heatmapPanel
}
