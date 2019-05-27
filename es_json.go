/**********************************************************
 * Author        : biliucheng
 * Email         : bilc_dev@163.com
 * Last modified : 2019-05-22 14:58
 * Filename      : es_json.go
 * Description   : json for generating dashboard struct.
 * Modified By   :
 * *******************************************************/
package autografana

var es_grafana_json string = `
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": null,
  "iteration": 1557993520350,
  "links": [],
  "panels": [
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 0,
        "y": 0
      },
      "id": 4,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "bucketAggs": [
            {
              "field": "@timestamp",
              "id": "2",
              "settings": {
                "interval": "10s",
                "min_doc_count": 0,
                "trimEdges": 0
              },
              "type": "date_histogram"
            }
          ],
          "metrics": [
            {
              "field": "METRIC_qps",
              "id": "1",
              "meta": {},
              "settings": {},
              "type": "avg"
            }
          ],
          "query": "FILTER_region:$FILTER_region  AND FILTER_az:$FILTER_az",
          "refId": "A",
          "timeField": "@timestamp"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "METRIC_qps",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 7,
        "y": 0
      },
      "id": 2,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "bucketAggs": [
            {
              "field": "@timestamp",
              "id": "2",
              "settings": {
                "interval": "10s",
                "min_doc_count": 0,
                "trimEdges": 0
              },
              "type": "date_histogram"
            }
          ],
          "metrics": [
            {
              "field": "METRIC_num",
              "id": "1",
              "meta": {},
              "settings": {},
              "type": "avg"
            }
          ],
          "query": "FILTER_region:$FILTER_region  AND FILTER_az:$FILTER_az",
          "refId": "A",
          "timeField": "@timestamp"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "METRIC_num",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 15,
        "y": 0
      },
      "id": 6,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "bucketAggs": [
            {
              "field": "@timestamp",
              "id": "2",
              "settings": {
                "interval": "10s",
                "min_doc_count": 0,
                "trimEdges": 0
              },
              "type": "date_histogram"
            }
          ],
          "metrics": [
            {
              "field": "select field",
              "id": "1",
              "type": "count"
            }
          ],
          "query": "FILTER_region:$FILTER_region  AND FILTER_az:$FILTER_az",
          "refId": "A",
          "timeField": "@timestamp"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "count",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    }
  ],
  "schemaVersion": 18,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": [
      {
        "allValue": null,
        "current": {
          "tags": [],
          "text": "All",
          "value": [
            "$__all"
          ]
        },
        "datasource": "test1",
        "definition": "{\"find\":\"terms\",\"field\":\"FILTER_region\"}",
        "hide": 0,
        "includeAll": true,
        "label": "FILTER_region",
        "multi": true,
        "name": "FILTER_region",
        "options": [],
        "query": "{\"find\":\"terms\",\"field\":\"FILTER_region\"}",
        "refresh": 2,
        "regex": "",
        "skipUrlSync": false,
        "sort": 0,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      },
      {
        "allValue": null,
        "current": {
          "tags": [],
          "text": "beijing",
          "value": [
            "beijing"
          ]
        },
        "datasource": "test1",
        "definition": "{\"find\":\"terms\",\"field\":\"FILTER_az\",\"query\":\"FILTER_region:$FILTER_region\"}",
        "hide": 0,
        "includeAll": true,
        "label": "FILTER_az",
        "multi": true,
        "name": "FILTER_az",
        "options": [],
        "query": "{\"find\":\"terms\",\"field\":\"FILTER_az\",\"query\":\"FILTER_region:$FILTER_region\"}",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 0,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      }
    ]
  },
  "time": {
    "from": "now/d",
    "to": "now/d"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "",
  "title": "test",
  "uid": null,
  "version": 8
}
`

var graph_panel_json string = `
{
  "aliasColors": {},
  "bars": false,
  "dashLength": 10,
  "dashes": false,
  "datasource": "grafana--smoke-qps*",
  "editable": true,
  "error": false,
  "fill": 1,
  "gridPos": {
    "h": 6,
    "w": 8,
    "x": 8,
    "y": 8
  },
  "id": 2,
  "isNew": false,
  "legend": {
    "alignAsTable": false,
    "avg": false,
    "current": false,
    "hideEmpty": false,
    "hideZero": false,
    "max": false,
    "min": false,
    "rightSide": false,
    "show": true,
    "total": false,
    "values": false
  },
  "lines": true,
  "linewidth": 1,
  "nullPointMode": "null",
  "percentage": false,
  "pointradius": 2,
  "points": false,
  "renderer": "flot",
  "seriesOverrides": [],
  "spaceLength": 10,
  "span": 0,
  "stack": false,
  "steppedLine": false,
  "targets": [
    {
      "bucketAggs": [
        {
          "field": "@timestamp",
          "id": "2",
          "settings": {
            "interval": "10s",
            "min_doc_count": 0
          },
          "type": "date_histogram"
        }
      ],
      "metrics": [
        {
          "field": "METRIC_qps",
          "id": "1",
          "meta": {},
          "settings": {},
          "type": "avg"
        }
      ],
      "query": "FILTER_region:$FILTER_region AND FILTER_user:$FILTER_user",
      "refId": "A",
      "timeField": "@timestamp"
    }
  ],
  "thresholds": [],
  "timeFrom": null,
  "timeRegions": [],
  "timeShift": null,
  "title": "METRIC_qps",
  "tooltip": {
    "shared": true,
    "sort": 0,
    "value_type": "individual"
  },
  "type": "graph",
  "xaxis": {
    "buckets": null,
    "format": "",
    "logBase": 0,
    "mode": "time",
    "name": null,
    "show": true,
    "values": []
  },
  "yaxes": [
    {
      "format": "short",
      "logBase": 1,
      "show": true
    },
    {
      "format": "short",
      "logBase": 1,
      "show": true
    }
  ],
  "yaxis": {
    "align": false,
    "alignLevel": null
  }
}
`

var heatmap_panel_json string = `
{
  "cards": {
    "cardPadding": null,
    "cardRound": null
  },
  "color": {
    "cardColor": "#56A64B",
    "colorScale": "sqrt",
    "colorScheme": "interpolateOranges",
    "exponent": 0.5,
    "mode": "opacity"
  },
  "dataFormat": "timeseries",
  "gridPos": {
    "h": 9,
    "w": 8,
    "x": 0,
    "y": 8
  },
  "heatmap": {},
  "hideZeroBuckets": false,
  "highlightCards": true,
  "id": 6,
  "legend": {
    "show": false
  },
  "reverseYBuckets": false,
  "timeFrom": null,
  "timeShift": null,
  "title": "Panel Title",
  "tooltip": {
    "show": true,
    "showHistogram": false
  },
  "type": "heatmap",
  "xAxis": {
    "show": true
  },
  "xBucketNumber": null,
  "xBucketSize": null,
  "yAxis": {
    "decimals": null,
    "format": "short",
    "logBase": 1,
    "max": null,
    "min": null,
    "show": true,
    "splitFactor": null
  },
  "yBucketBound": "auto",
  "yBucketNumber": null,
  "yBucketSize": null,
  "datasource": "grafana--replicaset-create-qps*",
  "targets": [
    {
      "refId": "A",
      "metrics": [
        {
          "type": "avg",
          "id": "1",
          "field": "select field"
        }
      ],
      "bucketAggs": [
        {
          "type": "date_histogram",
          "id": "2",
          "settings": {
            "interval": "10s",
            "min_doc_count": 0,
            "trimEdges": 0
          },
          "field": "@timestamp"
        }
      ],
      "timeField": "@timestamp"
    }
  ]
}
`