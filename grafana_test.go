package grafana

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	token        = "Bearer eyJrIjoiS0EydVpPMXZoMzlWbm1iTWo5dTZZaG1rWjhKMFRyeDUiLCJuIjoibGlkby1kaXNjb3JkLWJvdCIsImlkIjoxfQ=="
	addr         = "135.181.193.210:3001"
	dashboardUID = "xxixR7Z7z"

	timeout = 5 * time.Second
)

func TestPanels(t *testing.T) {
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
