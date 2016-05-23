package models

import (
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
)

type AlertRule struct {
	Id    int64 `json:"id"`
	OrgId int64 `json:"-"`
	//DataSourceId int64  `json:"datasourceId"`
	DashboardId  int64  `json:"dashboardId"`
	PanelId      int64  `json:"panelId"`
	Query        string `json:"query"`
	QueryRefId   string `json:"queryRefId"`
	WarnLevel    int64  `json:"warnLevel"`
	CritLevel    int64  `json:"critLevel"`
	WarnOperator string `json:"warnOperator"`
	CritOperator string `json:"critOperator"`
	Interval     string `json:"interval"`
	Frequency    int64  `json:"frequency"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	QueryRange   string `json:"queryRange"`
	Aggregator   string `json:"aggregator"`
	State        string `json:"state"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type AlertingClusterInfo struct {
	ServerId       string
	ClusterSize    int
	UptimePosition int
}

type HeartBeatCommand struct {
	ServerId string

	Result AlertingClusterInfo
}

type AlertRuleChange struct {
	Id      int64     `json:"id"`
	OrgId   int64     `json:"-"`
	AlertId int64     `json:"alertId"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
}

func (cmd *SaveDashboardCommand) GetAlertModels() []AlertRule {
	alerts := make([]AlertRule, 0)

	for _, rowObj := range cmd.Dashboard.Get("rows").MustArray() {
		row := simplejson.NewFromAny(rowObj)

		for _, panelObj := range row.Get("panels").MustArray() {
			panel := simplejson.NewFromAny(panelObj)

			alerting := panel.Get("alerting")
			alert := AlertRule{
				DashboardId:  cmd.Result.Id,
				OrgId:        cmd.Result.OrgId,
				PanelId:      panel.Get("id").MustInt64(),
				Id:           alerting.Get("id").MustInt64(),
				QueryRefId:   alerting.Get("queryRef").MustString(),
				WarnLevel:    alerting.Get("warnLevel").MustInt64(),
				CritLevel:    alerting.Get("critLevel").MustInt64(),
				WarnOperator: alerting.Get("warnOperator").MustString(),
				CritOperator: alerting.Get("critOperator").MustString(),
				Interval:     alerting.Get("interval").MustString(),
				Title:        alerting.Get("title").MustString(),
				Description:  alerting.Get("description").MustString(),
				QueryRange:   alerting.Get("queryRange").MustString(),
				Aggregator:   alerting.Get("aggregator").MustString(),
			}

			for _, targetsObj := range panel.Get("targets").MustArray() {
				target := simplejson.NewFromAny(targetsObj)

				if target.Get("refId").MustString() == alert.QueryRefId {
					targetJson, err := target.MarshalJSON()
					if err == nil {
						alert.Query = string(targetJson)
					}
					continue
				}
			}

			if alert.Query != "" {
				alerts = append(alerts, alert)
			}
		}
	}

	return alerts
}

// Commands
type SaveAlertsCommand struct {
	DashboardId int64
	UserId      int64
	OrgId       int64

	Alerts []AlertRule
}

type DeleteAlertCommand struct {
	AlertId int64
}

//Queries
type GetAlertsQuery struct {
	OrgId       int64
	State       []string
	DashboardId int64
	PanelId     int64

	Result []AlertRule
}

type GetAlertsForExecutionQuery struct {
	Timestamp int64

	Result []AlertRule
}

type GetAlertByIdQuery struct {
	Id int64

	Result AlertRule
}

type GetAlertChangesQuery struct {
	OrgId   int64
	Limit   int64
	SinceId int64

	Result []AlertRuleChange
}
