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

func TestPanelsFiltered(t *testing.T) {
	inst := NewGrafana(addr, token, timeout, ImageAttributes{
		Height:   500,
		Width:    1000,
		Timezone: "Europe/Moscow",
	})
	panels, err := inst.Panels(context.Background(), dashboardUID, "Slashing: Jailed Validators")
	if err != nil {
		t.Fatal(err)
	}
	if len(panels) != 1 {
		t.Fatal("List is not filtered")
	}
}

func TestGetPanelPicture(t *testing.T) {
	inst := NewGrafana(addr, token, timeout, ImageAttributes{
		Height:   500,
		Width:    1000,
		Timezone: "Europe/Moscow",
	})

	panels, _ := inst.Panels(context.Background(), dashboardUID)
	image := panels[0].Image
	imageBody, err := inst.GetPanelPicture(image)
	if err != nil {
		t.Fatal(err)
	}
	if len(imageBody) == 0 {
		t.Fatal("image is empty")
	}
}

func TestGetGrafanaPanel(t *testing.T) {
	inst := NewGrafana(addr, token, timeout, ImageAttributes{
		Height:   500,
		Width:    1000,
		Timezone: "Europe/Moscow",
	})

	panel, err := inst.GetGrafanaPanel("Slashing: Jailed Validators", dashboardUID)
	if err != nil {
		t.Fatal(err)
	}
	if panel.Title != "Slashing: Jailed Validators" {
		t.Fatal("Panel not found")
	}
}
