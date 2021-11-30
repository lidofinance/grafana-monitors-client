package grafana

import (
	"context"
	"fmt"
	"time"
)

const (
	imageURLFormat = "%s/render/d-solo/%s/lido-monitors?from=%d&to=%d&panelId=%d&width=%d&height=%d&tz=%s"

	httpPrefix = "http://"
)

type Grafana interface {
	Panels(ctx context.Context, dashboardUid string) ([]Panel, error)
}

type grafana struct {
	client *client
	attrs  ImageAttributes
}

func NewGrafana(url string, token string, timeout time.Duration, attrs ImageAttributes) Grafana {
	return &grafana{
		client: newClient(url, token, timeout),
		attrs:  attrs,
	}
}

func (g *grafana) Panels(ctx context.Context, dashboardUID string) ([]Panel, error) {
	var result = make(panelMap)

	dashboard, err := g.client.getDashboard(ctx, dashboardUID)
	if err != nil {
		return nil, fmt.Errorf("error getting dashboard response: %w", err)
	}

	alertStates, err := g.client.alertStates(ctx, dashboard.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting alert response: %w", err)
	}

	for _, p := range dashboard.Panels {
		currentValues, err := g.client.currentValues(ctx, p.Exprs)
		if err != nil {
			return nil, fmt.Errorf("error getting current values response: %w", err)
		}

		result[p.ID] = Panel{
			Title:         p.Title,
			CurrentValues: currentValues,
			Image:         g.getImageURL(dashboardUID, p.ID),
			Alert:         &p.Alert,
		}

		if as, ok := alertStates[p.ID]; ok {
			alert := result[p.ID].Alert

			alert.State = as.State
			alert.Name = as.Name
		}
	}

	return result.ToSlice(), nil

}

func (g *grafana) getImageURL(dashboardUID string, panelID int) string {
	to := time.Now()
	from := to.Add(-12 * time.Hour)

	return fmt.Sprintf(
		imageURLFormat,
		g.client.url,
		dashboardUID,
		from.UnixMilli(),
		to.UnixMilli(),
		panelID,
		g.attrs.Width,
		g.attrs.Height,
		g.attrs.Timezone,
	)
}
