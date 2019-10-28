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
	"strconv"

	sdk "github.com/bilc/grafana-sdk"
	"github.com/olivere/elastic"
	//elastic "gopkg.in/olivere/elastic.v6"
)

const PANEL_GRAPH = "graph"
const PAENL_HEATMAP = "heatmap"
const METRIC_INTERVAL= "interval"

type MyPanel struct {
	Title string `json:"title"`
	MyMetrics []MyMetric `json:"myMetrics"`
	Type string `json:"type"` // panel type: graph, heatmap
	Interval string `json:"interval"`
}

type MyMetric struct {
	Field string `json:"field"`
	Type string `json:"type"` //metric type: sum, avg, count...
}

func Es2Grafana(esUrl, esUrlNoAuth, esUser, esPassword, service, model string, grafanaUrl string, grafanaApiKey string,
	gratags, tagsSorts []string, tagsCascade map[string][]string, panels []MyPanel) (string, error ){

	ExpectTagsSort = tagsSorts
	index := IndexNameCommon(service, model)
	tags, metrics, err := ExtractEs(esUrl, index)
	if err != nil {
		return "", fmt.Errorf("extract es error: %s", err)
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
				return "", fmt.Errorf("tag %v not exist", j)
			}
		}
		tags = gratags
	}

	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)

	folderUid, err := GetFolderUid(grafanaUrl, grafanaApiKey, service)
	if err != nil{
		return "", fmt.Errorf("get folder uid err %v", err)
	}

	folderResp, err := grafanaCli.GetFolder(folderUid)
	if err != nil || folderResp.ID == 0 {
		folderResp, err = grafanaCli.CreateFolder(sdk.Folder{UID: FolderUid(service), Title: FolderUid(service)})
	}
	fmt.Println("Folder:", folderResp, err)

	ds := NewEsDataSource(esUrlNoAuth, index, esUser, esPassword)
	status, err := grafanaCli.CreateDatasource(ds)
	if err != nil {
		return "", fmt.Errorf("createdatasource err %v", err)
	}
	b, _ := json.Marshal(status)
	fmt.Println("---datasource: ", string(b))
	dashboard := NewGraphBoard(index, tags, tagsCascade, metrics, panels, model)
	b, _ = json.Marshal(dashboard)
	fmt.Println("---dashboard: ", string(b))

	resp, err := grafanaCli.SetDashboard(*dashboard, true, folderResp.ID)
	b, _ = json.Marshal(resp)
	fmt.Println("Debug dashboard rsp:", string(b))
	if err != nil {
		return "", fmt.Errorf("setDashboard err %v", err)
	}

	url, err := GetDashboardUrl(grafanaUrl, grafanaApiKey, service, model)
	if err != nil{
		return "", fmt.Errorf("getDashboardUrlByFolder err %s", err)
	}

	return url, nil
}

func GetDashboardUrl(grafanaUrl, grafanaApiKey, service, model string) (string,error) {
	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)
	folderUid, err := GetFolderUid(grafanaUrl, grafanaApiKey, service)
	if err != nil{
		return "", fmt.Errorf("get folder uid err %v", err)
	}

	folderResp, err := grafanaCli.GetFolder(folderUid)
	if err != nil{
		return "", fmt.Errorf("get folders err %v", err)
	}

	folders,err := grafanaCli.SearchFolders(folderResp.ID)
	if err != nil{
		return "", fmt.Errorf("search folders err %v", err)
	}
	for _, dashboard := range folders{
		if dashboard.Title == model{
			return dashboard.URL, nil
		}
	}
	return "", nil
}

func GetAllDashboardInFolder(grafanaUrl, grafanaApiKey, service string) ([]sdk.FoundFolder, error){
	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)
	folderUid, err := GetFolderUid(grafanaUrl, grafanaApiKey, service)
	if err != nil || folderUid == ""{
		return nil, fmt.Errorf("get folder uid err[%v] or folder[%s] not exist ", err, service)
	}
	folderResp, err := grafanaCli.GetFolder(folderUid)
	if err != nil{
		return nil, fmt.Errorf("get folders err %v", err)
	}

	folders,err := grafanaCli.SearchFolders(folderResp.ID)
	if err != nil{
		return nil, fmt.Errorf("search folders err %v", err)
	}
	return folders, nil
}

func GetFolderUid(grafanaUrl, grafanaApiKey, folderTitle string) (string, error){
	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)
	allFolders, err := grafanaCli.SearchFolders(0)
	if err != nil{
		return "", fmt.Errorf("list folders err %v", err)
	}

	for _ , folder := range allFolders{
		if folder.Title == folderTitle{
			return folder.UID, nil
		}
	}
	return "", nil
}

func ListFolders(grafanaUrl, grafanaApiKey string) ([]sdk.FoundFolder, error){
	grafanaCli := sdk.NewClient(grafanaUrl, grafanaApiKey, sdk.DefaultHTTPClient)
	allFolders, err := grafanaCli.SearchFolders(0)
	if err != nil{
		return nil, fmt.Errorf("list folders err %v", err)
	}
	return allFolders, err
}

func ListServiceModel(esUrl, index string) ([]string, error) {
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
		return nil, fmt.Errorf("get es index error %v", err)
	}
	models := make([]string, 0)
	prefix := strings.TrimRight(index, "*") + "-"
	for replyKey := range reply {
		if strings.HasPrefix(replyKey, prefix) {
			modelTail := strings.TrimPrefix(replyKey, prefix)
			lastIndex := strings.LastIndex(modelTail, "-")
			model := modelTail[0:lastIndex]
			models = append(models, model)
		}
	}

	fmt.Println("model ", RemoveRepeatedElement(models))
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
		return nil, nil, fmt.Errorf("get es insdex err %v", err)
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
		for field := range tmp {
			if strings.HasPrefix(field, FIELD_TAG_PREFIX) {
				tags = append(tags, field)
			} else if strings.HasPrefix(field, FIELD_METRIC_PREFIX) || strings.HasPrefix(field, FIELD_SUM_METRIC_PREFIX)  {
				metrics = append(metrics, field)
			}
		}
	}
	sort.Strings(tags)
	sort.Strings(metrics)
	fmt.Println(tags, metrics)
	return tags, metrics, nil
}

func NewEsDataSource(esUrl string, db string, user,password string) sdk.Datasource {
	//	jsonData
	//	esVersion: 60
	//keepCookies: []
	//timeField: "@timestamp"
	tmp := true
	ds := sdk.Datasource{
		Access:    "proxy",
		BasicAuth: &tmp,
		BasicAuthUser: &user,
		BasicAuthPassword: &password,
		Name:      db,
		Database:  &db,
		URL:       esUrl,
		Type:      "elasticsearch",
		JSONData:  map[string]interface{}{"esVersion": "60", "timeField": "@timestamp", "keepCookies": []string{}},
	}
	return ds
}

func NewGraphBoard(myDataSource string, mytags []string, mytagsCascade map[string][]string, myMetrics []string, myPanels []MyPanel, myTitle string) *sdk.Board {
	var myID uint = 1
	var board sdk.Board
	err := json.Unmarshal([]byte(es_grafana_json), &board)
	if err != nil {
		fmt.Println("grafana json unmarshal error", err)
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
		// if tag in mytagsCascade, then setting query cascaded
		if relation, ok := mytagsCascade[tag]; ok {
			queryPrefix := fmt.Sprintf("{\"find\":\"terms\",\"field\":\"%s\",\"query\":\"", tag)
			query := ""
			querySuffix := "\"}"
			existTags := getExistTags(relation, mytags)
			for i, r := range existTags {
				if i == len(existTags)-1 {
					query += fmt.Sprintf("%s:$%s", r, r)
				} else {
					query += fmt.Sprintf("%s:$%s AND ", r, r)
				}
			}
			templateVar.Query = queryPrefix + query + querySuffix
		} else {
			templateVar.Query = fmt.Sprintf("{\"find\":\"terms\",\"field\":\"%s\"}", tag)
		}
		//	templateVar.Definition = templateVar.Query
		board.Templating.List = append(board.Templating.List, templateVar)
	}
	SortTemplatingList(board.Templating.List)

	if len(luceneQuery) > 5 {
		luceneQuery = luceneQuery[0 : len(luceneQuery)-5]
	}
	//这里使用指针，如果*board.Panels[0]会丢失数据
	//panelb, _ := json.Marshal(board.Panels[0])
	//	panelVar := *board.Panels[0]
	board.Panels = board.Panels[0:0]
	for index, myPanel :=  range myPanels{
		var panel *sdk.Panel
		if myPanel.Type == PANEL_GRAPH{
			panel = NewGraphPanel(myDataSource, myID, index, myPanel.Title, myPanel.Interval, myPanel.MyMetrics, luceneQuery)
		}else if myPanel.Type == PAENL_HEATMAP{
			panel = NewHeatmapPanel(myDataSource, myID,index, myPanel.Title, myPanel.Interval, myPanel.MyMetrics, luceneQuery)
		}
		myID += 1
		board.Panels = append(board.Panels, panel)
	}
	return &board
}

type TemplateVars []sdk.TemplateVar

func getIndex(name string, arrays []string) int {
	for index, arr := range arrays {
		if strings.EqualFold(arr, name) {
			return index
		}
	}
	return -1
}

func getExistTags(tags, mytags []string) []string {
	var existTag []string
	for _, tag := range tags {
		for _, mytag := range mytags {
			if tag == mytag {
				existTag = append(existTag, tag)
			}
		}
	}
	return existTag
}

func isExist(name string, arrays []string) bool {
	for _, arr := range arrays {
		if arr == name {
			return true
		}
	}
	return false
}

// sort TAG_* expected and in-situ output user-defined
func (c TemplateVars) Less(i, j int) bool {
	//c[i].Name < c[j].Name
	indexI := getIndex(c[i].Name, ExpectTagsSort)
	indexJ := getIndex(c[j].Name, ExpectTagsSort)

	if indexI != -1 && indexJ != -1 {
		return indexI < indexJ
	} else {
		return indexI > indexJ
	}
}

func (c TemplateVars) Len() int {
	return len(c)
}

func (c TemplateVars) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func SortTemplatingList(templateVars TemplateVars) {
	sort.Sort(templateVars)
}

func FolderUid(service string) string {
	return service
}

func NewGraphPanel(myDataSource string, panelId uint, panelIndex int, panelTile, panelInterval string, myMetrics []MyMetric, luceneQuery string) *sdk.Panel {
	var graphPanel sdk.Panel
	err := json.Unmarshal([]byte(graph_panel_json), &graphPanel)
	if err != nil {
		fmt.Println("unmarshl graph panel json error:", err)
		return nil
	}

	graphPanel.Datasource = &myDataSource
	*graphPanel.GridPos.X = (panelIndex % 3) * 8
	*graphPanel.GridPos.Y = (panelIndex / 3) * 8
	graphPanel.ID = panelId
	graphPanel.Title = panelTile

	graphPanel.GraphPanel.Targets[0].Query = luceneQuery
	graphPanel.GraphPanel.Targets[0].BucketAggs[0].Settings.Interval = panelInterval

	metrics := make([]sdk.Metric, len(myMetrics))
	for i, metric := range myMetrics{
		metrics[i].ID = strconv.Itoa(i)
		metrics[i].Field = metric.Field
		metrics[i].Type = metric.Type
		metrics[i].Meta = struct {}{}
		metrics[i].Settings = struct {}{}
	}
	graphPanel.GraphPanel.Targets[0].Metrics = metrics
	return &graphPanel
}

func NewHeatmapPanel(myDataSource string, panelId uint, panelIndex int, panelTile, panelInterval string,  myMetrics []MyMetric, luceneQuery string) *sdk.Panel {
	var heatmapPanel sdk.Panel
	err := json.Unmarshal([]byte(heatmap_panel_json), &heatmapPanel)
	if err != nil {
		fmt.Println("unmarshl heatmap panel json error:", err)
		return nil
	}

	heatmapPanel.Datasource = &myDataSource
	*heatmapPanel.GridPos.X = (panelIndex % 3) * 8
	*heatmapPanel.GridPos.Y = (panelIndex / 3) * 8
	heatmapPanel.ID = panelId
	heatmapPanel.Title = panelTile

	heatmapPanel.HeatmapPanel.Targets[0].Query = luceneQuery
	heatmapPanel.HeatmapPanel.Targets[0].BucketAggs[0].Settings.Interval = panelInterval

	metrics := make([]sdk.Metric, len(myMetrics))
	for i, metric := range myMetrics{
		metrics[i].ID = strconv.Itoa(i)
		metrics[i].Field = metric.Field
		metrics[i].Type = metric.Type
		metrics[i].Meta = struct {}{}
		metrics[i].Settings = struct {}{}
	}
	heatmapPanel.HeatmapPanel.Targets[0].Metrics = metrics
	return &heatmapPanel
}

/**
 use this method to get metric interval
 when metric="interval-20s_METRIC_cpu_util" return "20s"
 when metric="METRIC_cpu_util" return "10s"
*/
var getMetricInterval = func(metric string) (interval string){
	if strings.HasPrefix(metric, METRIC_INTERVAL){
		interval = strings.TrimLeft(strings.Split(metric, "_")[0], METRIC_INTERVAL+ "-")
	}else{
		interval = "10s"
	}
	return
}

/**
 use this method to get metric field
 when metric="interval-20s_METRIC_cpu_util" return "METRIC_cpu_util"
 when metric="METRIC_cpu_util" return "METRIC_cpu_util"
*/
var getMetricField = func(metric string) (field string){
	if strings.HasPrefix(metric, METRIC_INTERVAL){
		field = strings.TrimLeft(metric, strings.Split(metric, "_")[0] + "_")
	}else{
		field = metric
	}
	return
}

/**
 use thie method to get metric type
 when metric="interval-20s_METRIC_cpu_util" return "avg"
 when metric="SUM_METRIC_cpu_util" return "sum"
*/
var getMetricType = func(metric string) (mType string){
	if strings.HasPrefix(getMetricField(metric), FIELD_SUM_METRIC_PREFIX){
		mType = "sum"
	}else{
		mType = "avg"
	}
	return
}


