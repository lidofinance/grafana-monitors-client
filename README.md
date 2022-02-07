# Grafana monitors client

To run tests please make sure you provided following environment variables

- `GRAFANA_TOKEN`
- `GRAFANA_ADDR`
- `GRAFANA_DASHBOARD_UID`

Example:

```
GRAFANA_TOKEN="<Your Grafana token>" GRAFANA_ADDR=<your Grafana address> GRAFANA_DASHBOARD_UID=<your Grafana dashboard Id> go test

#or for short tests

GRAFANA_TOKEN="<Your Grafana token>" GRAFANA_ADDR=<your Grafana address> GRAFANA_DASHBOARD_UID=<your Grafana dashboard Id> go test -short

```
