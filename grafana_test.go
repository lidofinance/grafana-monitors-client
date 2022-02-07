package grafana

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	timeout = 5 * time.Second
)

var token string
var addr string
var dashboardUID string

func TestPanels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	dashboardUID = os.Getenv("GRAFANA_DASHBOARD_UID")
	if dashboardUID == "" {
		t.Fatal("Failed grafana test: no dashboard uid")
	}
	inst := NewGrafana(addr, token, timeout, ImageAttributes{
		Height:   500,
		Width:    1000,
		Timezone: "Europe/Moscow",
	})

	panels, err := inst.Panels(context.Background(), dashboardUID)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range panels {
		if p.Title == "" {
			t.Fatal("title is empty!")
		}

		if len(p.CurrentValues) == 0 {
			t.Fatalf("current values is empty! title is %s", p.Title)
		}

		resp, err := http.Get(p.Image)
		if err != nil {
			t.Fatalf("failed to get image %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read image body: %s", err)
		}

		if len(body) == 0 {
			t.Fatalf("image is empty: %s", err)
		}
	}
}
