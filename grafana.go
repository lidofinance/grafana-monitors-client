package grafana

import (
	"context"
	"fmt"
	"strings"
	"sync"
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

	panelsCurrentValues, err := g.getPanelsCurrentValues(ctx, dashboard)
	if err != nil {
		return nil, fmt.Errorf("faild to getPanelsCurrentValues: %w", err)
	}

	for p, currentValues := range panelsCurrentValues {
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

func (g *grafana) getPanelsCurrentValues(ctx context.Context, dashboard dashboardData) (map[*panelData][]CurrentValue, error) {
	var (
		wg = &sync.WaitGroup{}
		errsCh = make(chan error, len(dashboard.Panels))
		panelsCurrentValues = map[*panelData][]CurrentValue{}
	)
	for _, p := range dashboard.Panels {
		go func (p *panelData, panelsCurrentValues map[*panelData][]CurrentValue) {
			wg.Add(1)
			defer wg.Done()

			currentValues, err := g.client.currentValues(ctx, p.Exprs)
			if err != nil {
				errsCh <- fmt.Errorf("failed to get currentValues for panel (%s): %w", p.Title, err)
				return
			}

			panelsCurrentValues[p] = currentValues
		}(&p, panelsCurrentValues)
	}
	wg.Wait()

	// Close the errors channel so that when we hit the end of the errors queue, the
	// range loop stops.
	close(errsCh)
	var errs []string
	for workerErr := range errsCh {
		errs = append(errs, workerErr.Error())
	}
	if len(errs) != 0 {
		return nil, fmt.Errorf("collected errors: %s", strings.Join(errs, ";"))
	}

	return panelsCurrentValues, nil
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
