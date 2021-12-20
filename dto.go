package grafana

const (
	rowPanelType = "row"
)

type alertStateDTO struct {
	ID            int    `json:"id"`
	DashboardID   int    `json:"dashboardId"`
	DashboardUID  string `json:"dashboardUid"`
	DashboardSlug string `json:"dashboardSlug"`
	PanelID       int    `json:"panelId"`
	Name          string `json:"name"`
	State         string `json:"state"`
	NewStateDate  string `json:"newStateDate"`
}

type alertStatesDTO []alertStateDTO

func (as alertStatesDTO) ToAlertMap() map[int]Alert {
	alertsMap := make(map[int]Alert)

	for _, a := range as {
		alertsMap[a.PanelID] = Alert{
			Name:  a.Name,
			State: a.State,
		}
	}

	return alertsMap
}

type dashboardDTO struct {
	Dashboard struct {
		ID     int     `json:"id"`
		Panels []panel `json:"panels"`
		UID    string  `json:"uid"`
	} `json:"dashboard"`
}

type panel struct {
	ID    int `json:"id"`
	Alert struct {
		Conditions []struct {
			Evaluator struct {
				Params []float64 `json:"params"`
				Type   string    `json:"type"`
			} `json:"evaluator"`
			Operator struct {
				Type string `json:"type"`
			} `json:"operator"`
			Query struct {
				Params []string `json:"params"`
			} `json:"query"`
			Reducer struct {
				Params []interface{} `json:"params"`
				Type   string        `json:"type"`
			} `json:"reducer"`
			Type string `json:"type"`
		} `json:"conditions"`
		ExecutionErrorState string        `json:"executionErrorState"`
		For                 string        `json:"for"`
		Frequency           string        `json:"frequency"`
		Handler             int64         `json:"handler"`
		Message             string        `json:"message"`
		Name                string        `json:"name"`
		NoDataState         string        `json:"noDataState"`
		Notifications       []interface{} `json:"notifications"`
	} `json:"alert"`
	Targets []struct {
		Expr         string `json:"expr"`
		LegendFormat string `json:"legendFormat"`
	} `json:"targets"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

func (d *dashboardDTO) Data() (result dashboardData) {
	result.ID = d.Dashboard.ID

	for _, p := range d.Dashboard.Panels {
		if p.Type == rowPanelType {
			continue
		}

		panel := panelData{
			ID:    p.ID,
			Title: p.Title,
		}

		for _, t := range p.Targets {
			panel.Exprs = append(panel.Exprs, expr{
				Query:        t.Expr,
				LegendFormat: t.LegendFormat,
			})
		}

		panel.Alert = Alert{
			Name: p.Alert.Name,
		}

		for _, c := range p.Alert.Conditions {
			panel.Alert.Conditions = append(panel.Alert.Conditions, Condition{
				Type:   c.Evaluator.Type,
				Values: c.Evaluator.Params,
			})

		}

		result.Panels = append(result.Panels, panel)
	}

	return result
}

type datasourceDTO struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Label    string `json:"label"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
	Error string `json:"error"`
}

func (d datasourceDTO) ToLabelValues() []LabelValue {
	var currentLabelValues []LabelValue

	for _, r := range d.Data.Result {
		if len(r.Values) == 0 {
			return nil
		}

		values := r.Values[0]

		// first element in array is unix timestamp, second element is current value in string type
		if len(values) == 2 {
			value, ok := values[1].(string)
			if !ok {
				return nil
			}

			currentLabelValues = append(currentLabelValues, LabelValue{
				Label: r.Metric.Label,
				Value: value,
			})
		}
	}

	return currentLabelValues
}
