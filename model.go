package grafana

type Panel struct {
	Title         string         `json:"title"`
	Image         string         `json:"image"`
	Alert         *Alert         `json:"alert"`
	CurrentValues []CurrentValue `json:"current_value"`
}

type Alert struct {
	Name       string      `json:"name"`
	State      string      `json:"state"`
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Type   string    `json:"type"` // gt, lt, etc...
	Values []float64 `json:"values"`
}

type CurrentValue struct {
	Query  string       `json:"query"`
	Values []LabelValue `json:"values"`
}

type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type panelMap map[int]Panel

func (p panelMap) ToSlice() []Panel {
	var result []Panel

	for _, panel := range p {
		result = append(result, panel)
	}

	return result
}

type dashboardData struct {
	ID     int
	Panels []panelData
}

type panelData struct {
	ID    int
	Title string
	Exprs []expr
	Alert Alert
}

type expr struct {
	Query        string
	LegendFormat string
}

type ImageAttributes struct {
	Width    int
	Height   int
	Timezone string
}
