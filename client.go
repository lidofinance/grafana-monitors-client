package grafana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dashboardPath   = "/api/dashboards/uid/"
	alertsPath      = "/api/alerts"
	datasourcesPath = "/api/datasources/proxy/1/api/v1/query_range"
)

const (
	authHeader = "Authorization"

	// in our case it doesn't matter, but required
	defaultStepQueryParam = 10
	multipleLegendFormat  = "{{label}}"
)

type client struct {
	url    string
	token  string
	client http.Client
}

func newClient(url string, token string, timeout time.Duration) *client {
	if strings.Index(url, httpPrefix) != 0 {
		url = httpPrefix + url
	}

	return &client{
		client: http.Client{Timeout: timeout},
		url:    url,
		token:  token,
	}
}

func (c *client) getDashboard(ctx context.Context, dashboardUID string) (dashboardData, error) {
	var dashboard dashboardDTO

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s%s", c.url, dashboardPath, dashboardUID),
		nil)
	if err != nil {
		return dashboardData{}, fmt.Errorf("failed to create NewRequestWithContext: %w", err)
	}

	req.Header.Add(authHeader, c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return dashboardData{}, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dashboardData{}, fmt.Errorf("failed: status code is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dashboardData{}, fmt.Errorf("failed to read dashboard response body: %w", err)
	}

	if err = json.Unmarshal(body, &dashboard); err != nil {
		return dashboardData{}, fmt.Errorf("failed to unmarshal dashboard response: %w", err)
	}

	return dashboard.Data(), nil
}

func (c *client) alertStates(ctx context.Context, dashboardID int) (map[int]Alert, error) {
	var alertStates alertStatesDTO

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", c.url, alertsPath),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewRequestWithContext: %w", err)
	}

	req.URL.Query().Add("dashboardId", strconv.FormatInt(int64(dashboardID), 10))
	req.Header.Add(authHeader, c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed: status code is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read alert states response body: %w", err)
	}

	if err = json.Unmarshal(body, &alertStates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert states response: %w", err)
	}

	return alertStates.ToAlertMap(), nil
}

func (c *client) currentValues(ctx context.Context, queries []expr) ([]CurrentValue, error) {
	result := make([]CurrentValue, 0, len(queries))

	for _, query := range queries {
		currentLabelValues, err := c.datasource(ctx, query.Query)
		if err != nil {
			return nil, fmt.Errorf("error getting current values by query: %s; error: %s", query, err)
		}

		if query.LegendFormat != multipleLegendFormat {
			for i := range currentLabelValues {
				currentLabelValues[i].Label = query.LegendFormat
			}
		}

		result = append(result, CurrentValue{
			Query:  query.Query,
			Values: currentLabelValues,
		})
	}

	return result, nil
}

func (c *client) datasource(ctx context.Context, query string) ([]LabelValue, error) {
	var datasource datasourceDTO

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", c.url, datasourcesPath),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewRequestWithContext: %w", err)
	}

	now := time.Now().Unix()

	q := req.URL.Query()

	q.Add("query", query)
	q.Add("start", strconv.FormatInt(now, 10))
	q.Add("end", strconv.FormatInt(now, 10))
	q.Add("step", strconv.FormatInt(defaultStepQueryParam, 10))

	req.URL.RawQuery = q.Encode()

	req.Header.Add(authHeader, c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read dashboard response body: %s: %w", body, err)
	}

	if err = json.Unmarshal(body, &datasource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard response: %w", err)
	}

	if datasource.Error != "" {
		return nil, errors.New(datasource.Error)
	}

	return datasource.ToLabelValues(), nil
}
