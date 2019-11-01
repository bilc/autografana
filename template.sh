#!/bin/bash

ESDOMAIN=${ESDOMAIN-localhost:9200}
curl -XPUT "http://${ESDOMAIN}/_template/grafana-template" -H 'Content-Type: application/json' -d'
{
    "index_patterns":["grafana--*"],
    "settings" : {
      "index" : {
        "number_of_shards" : "1",
        "number_of_replicas" : "1"
      }
    },
    "mappings": {
        "_doc":{
            "properties":{
                "@timestamp":{"type": "date"},
                "service":{"type": "keyword"},
                "model":{"type": "keyword"}
            },
            "dynamic_templates": [
                {
                    "match_tag":{
                        "match":"TAG_*",
                        "mapping": {
                            "type": "keyword"
                        }
                    }
                },
                {
                    "match_number":{
                        "match":"METRIC_*",
                        "mapping": {
                            "type": "long"
                        }
                    }
                },
                {
                    "match_number" : {
                        "match" : "SUM_METRIC_*",
                        "mapping" : {
                            "type" : "long"
                         }
                    }
                },
                {
                    "match_number" : {
                        "match" : "AVG-GRAPH_*",
                        "mapping" : {
                            "type" : "long"
                         }
                    }
                },
                {
                    "match_number" : {
                        "match" : "SUM-GRAPH_*",
                        "mapping" : {
                            "type" : "long"
                         }
                    }
                },
                {
                    "match_number" : {
                        "match" : "HEATMAP_*",
                        "mapping" : {
                            "type" : "long"
                         }
                    }
                },
                {
                    "match_search":{
                        "match":"SEARCH_*",
                        "mapping": {
                            "type": "text"
                        }
                    }
                },
                {
                    "strings_tpl": {
                        "match_mapping_type": "string",
                        "mapping": {
                            "type": "keyword"
                        }
                    }
                }
            ]
        }
    }
}'
