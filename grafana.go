package grafana

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	imageURLFormat = "%s/render/d-solo/%s/lido-monitors?from=%d&to=%d&panelId=%d&width=%d&height=%d&tz=%s"
	httpPrefix     = "http://"
)

type Grafana interface {
	Panels(ctx context.Context, dashboardUid string, filterPanelNames ...string) ([]Panel, error)
	GetPanelPicture(url string) ([]byte, error)
	GetGrafanaPanel(panelName string, dashboardID string) (*Panel, error)
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

func (g *grafana) Panels(ctx context.Context, dashboardUID string, filterPanelNames ...string) ([]Panel, error) {
	var result = make(panelMap)
	dashboard, err := g.client.getDashboard(ctx, dashboardUID)

	if err != nil {
		return nil, fmt.Errorf("error getting dashboard response: %w", err)
	}

	alertStates, err := g.client.alertStates(ctx, dashboard.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting alert response: %w", err)
	}

	panels := make([]panelData, 0)
	if len(filterPanelNames) > 0 {
		panelNamesMap := make(map[string]bool)
		for _, pn := range filterPanelNames {
			panelNamesMap[pn] = true
		}
		for _, panel := range dashboard.Panels {
			if _, ok := panelNamesMap[panel.Title]; ok {
				panels = append(panels, panel)
			}
		}
	} else {
		panels = dashboard.Panels
	}

	for _, p := range panels {
		currentValues, err := g.client.currentValues(ctx, p.Exprs)
		if err != nil {
			return nil, fmt.Errorf("error getting current values response: %w", err)
		}

		result[p.ID] = Panel{
			Title:         p.Title,
			CurrentValues: currentValues,
			Image:         g.getImageURL(dashboardUID, p.ID),
			Alert:         p.Alert,
		}

		if as, ok := alertStates[p.ID]; ok {
			panel := result[p.ID]

			panel.Alert.State = as.State
			panel.Alert.Name = as.Name

			result[p.ID] = panel
		}
	}

	return result.ToSlice(), nil

}

func (g *grafana) GetPanelPicture(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Add(authHeader, g.client.token)
	resp, err := g.client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get a panel picture: status code is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read a panel picture body: %w", err)
	}

	return body, nil
}

func (g *grafana) GetGrafanaPanel(panelName string, dashboardID string) (*Panel, error) {
	panels, err := g.Panels(context.Background(), dashboardID, panelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get grafana panels: %w", err)
	}
	if len(panels) == 1 {
		return &panels[0], nil
	}

	return nil, fmt.Errorf("panel with name %s not found", panelName)
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
