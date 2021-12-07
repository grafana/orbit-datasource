package orbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Figure struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			TimeFrameParams struct {
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
				GroupBy   string `json:"group_by"`
			} `json:"timeframe_params"`
			ViewData struct {
				Data []TableDataPoint `json:"data"`
			} `json:"view_data"`
			ViewParams struct {
				ViewType string `json:"view_type"`
			} `json:"view_params"`
		} `json:"attributes"`
	} `json:"data"`
}

type TableDataPoint struct {
	Name  string `json:"name"`
	Date  string `json:"x"`
	Value int64  `json:"y"`
}

type Client struct {
	httpClient *http.Client
	apiToken   string
	workspace  string
}

type ClientOption func(client *Client)

func NewClient(workspace string, apiToken string, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{},
		apiToken:   apiToken,
		workspace:  workspace,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

type FigureQuery struct {
	Analyze       string   `json:"analyze"`
	NewReturning  string   `json:"newReturning"`
	Interval      string   `json:"interval"`
	Affiliation   string   `json:"affiliation"`
	GroupBy       string   `json:"groupBy"`
	GroupLimit    string   `json:"groupLimit"`
	Orbits        []string `json:"orbits"`
	ActivityTypes []string `json:"activityTypes"`
	Cumulative    bool     `json:"cumulative"`
}

func (c *Client) Test(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://app.orbit.love/%s/figures/new.json", c.workspace), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("unexpected status code")
	}

	return 200, nil
}

func (c *Client) Figure(ctx context.Context, startDate, endDate string, figureQuery FigureQuery) (*Figure, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://app.orbit.love/%s/figures/new.json", c.workspace), nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("start_date", startDate)
	query.Set("end_date", endDate)
	query.Set("figure_data", figureQuery.Analyze)
	query.Set("group_by", figureQuery.Interval)

	if figureQuery.Cumulative {
		query.Set("cumulative", "true")
	}

	if len(figureQuery.Orbits) > 0 {
		query.Set("orbit", strings.Join(figureQuery.Orbits, ""))
	}
	if len(figureQuery.ActivityTypes) > 0 {
		query.Set("activity_type", url.PathEscape(strings.Join(figureQuery.ActivityTypes, ",")))
	}

	if figureQuery.NewReturning != "" {
		query.Set("new_returning", figureQuery.NewReturning)
	}

	if figureQuery.GroupBy != "" {
		query.Set("figure_group_by", figureQuery.GroupBy)

		if figureQuery.GroupLimit != "" {
			query.Set("figure_group_limit", figureQuery.GroupLimit)
		}
	}

	if figureQuery.Affiliation != "" {
		query.Set("affiliation", figureQuery.Affiliation)
	}

	query.Set("view_type", "table")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.DefaultLogger.Error(req.URL.String())
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var figure Figure
	if err := json.NewDecoder(resp.Body).Decode(&figure); err != nil {
		return nil, err
	}

	return &figure, nil
}

type ActivityType struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Attributes struct {
		Name      string `json:"name,omitempty"`
		ShortName string `json:"short_name,omitempty"`
		Key       string `json:"key,omitempty"`
		Channel   string `json:"channel,omitempty"`
		Weight    string `json:"weight,omitempty"`
	} `json:"attributes,omitempty"`
}

func (c *Client) ActivityTypes(ctx context.Context) ([]ActivityType, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://app.orbit.love/api/v1/%s/activity_types", c.workspace), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.DefaultLogger.Error(req.URL.String())
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payload struct {
		Data []ActivityType `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data, nil
}
