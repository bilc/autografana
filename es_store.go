/**********************************************************
 * Author        : biliucheng
 * Email         : bilc_dev@163.com
 * Last modified : 2019-05-22 12:59
 * Filename      : es_store.go
 * Description   : put document to elasticsearch
 * Modified By   :
 * *******************************************************/
package autografana

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

var INDEX_PREFIX string = "grafana--"
var ALIAS_PREFIX string = "graf-alias--"

var FIELD_SERVICE string = "service"
var FIELD_MODEL string = "model"
var FIELD_TIMESTAMP string = "@timestamp"
var FIELD_TAG_PREFIX string = "TAG_"
var FIELD_METRIC_PREFIX string = "METRIC_"

var FieldErr error = errors.New("field wrong")

type EsClient struct {
	*elastic.Client

	indices map[string]interface{}
}

func NewEsClient(url string) (*EsClient, error) {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &EsClient{client, make(map[string]interface{})}, nil
}

func (e *EsClient) PutDoc(doc map[string]interface{}) error {
	var service, model string
	if s, ok := doc[FIELD_SERVICE]; ok {
		service = s.(string)
	}
	if m, ok := doc[FIELD_MODEL]; ok {
		model = m.(string)
	}
	if service == "" || model == "" {
		return fmt.Errorf("%v, %v %v :%v", FieldErr, service, model, doc)
	}
	index := IndexName(service, model)

	_, err := e.Client.Index().Index(index).BodyJson(doc).Do(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "index_not_found_exception") {
			if _, err = e.Client.CreateIndex(index).Do(context.Background()); err == nil {
				_, err = e.Client.Index().Index(index).BodyJson(doc).Do(context.Background())
			}
		}
	}
	return err
}

func (e *EsClient) CreateIndex(index string) {
	e.Client.CreateIndex(index).Do(context.Background())
}

func IndexName(service, model string) string {
	n := time.Now()
	return fmt.Sprintf("%v-%v%02d", INDEX_PREFIX+service+"-"+model, n.Year(), n.Month())
}

func IndexNameCommon(service, model string) string {
	return INDEX_PREFIX + service + "-" + model + "*"
}

//func IndexAlias(service, model string) string {
//	return ALIAS_PREFIX + service + "-" + model
//}
