package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/orbit-datasource/pkg/orbit"
)

var (
	_ backend.QueryDataHandler    = (*OrbitDatasource)(nil)
	_ backend.CheckHealthHandler  = (*OrbitDatasource)(nil)
	_ backend.CallResourceHandler = (*OrbitDatasource)(nil)
)

// queryModel represents the query that's constructed by the query editor and
// sent by the dashboard.
type queryModel struct {
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

type OrbitDatasource struct {
	client *orbit.Client
}

// NewOrbitDatasource creates a new datasource instance.
func NewOrbitDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var jsonData struct {
		Workspace string `json:"workspacePath"`
	}

	if err := json.Unmarshal(settings.JSONData, &jsonData); err != nil {
		return nil, err
	}

	apiToken := settings.DecryptedSecureJSONData["apiToken"]

	ds := &OrbitDatasource{
		client: orbit.NewClient(jsonData.Workspace, apiToken),
	}

	return ds, nil
}

func (d *OrbitDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		response.Responses[q.RefID] = d.query(ctx, req.PluginContext, q)
	}
	return response, nil
}

func (d *OrbitDatasource) query(ctx context.Context, pluginCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	var qm queryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		response.Error = err
		return response
	}

	figQuery := orbit.FigureQuery{
		Analyze:       qm.Analyze,
		Affiliation:   qm.Affiliation,
		Interval:      qm.Interval,
		NewReturning:  qm.NewReturning,
		GroupBy:       qm.GroupBy,
		GroupLimit:    qm.GroupLimit,
		Orbits:        qm.Orbits,
		ActivityTypes: qm.ActivityTypes,
		Cumulative:    qm.Cumulative,
	}

	fig, err := d.client.Figure(ctx, query.TimeRange.From.Format("2006-01-02"), query.TimeRange.To.Format("2006-01-02"), figQuery)
	if err != nil {
		response.Error = err
		return response
	}

	frame, err := figureToWideFrame(fig)
	if err != nil {
		response.Error = err
		return response
	}

	response.Frames = data.Frames{frame}

	return response
}

// CheckHealth runs when the user presses "Save & Test" in the data source settings. Tests that the user can access
// their workspace.
func (d *OrbitDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	status, err := d.client.Test(ctx)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "Couldn't find a workspace with that name.",
			}, nil
		case http.StatusUnauthorized:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "You're not authorized to access the workspace.",
			}, nil
		default:
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: fmt.Sprintf("Received an unexpected status code: %d", status),
			}, nil
		}
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Success",
	}, nil
}

// CallResources exposes a REST API with support operations for the query editor.
func (d *OrbitDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	activities, err := d.client.ActivityTypes(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(activities)
	if err != nil {
		return err
	}

	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   b,
	})
}

// figureToWideFrame converts an Orbit figure to a Grafana data frame.
func figureToWideFrame(figure *orbit.Figure) (*data.Frame, error) {
	dates := map[string]bool{}
	series := map[string]map[string]int64{}

	for _, point := range figure.Data.Attributes.ViewData.Data {
		if _, ok := series[point.Name]; !ok {
			series[point.Name] = make(map[string]int64)
		}
		series[point.Name][point.Date] = point.Value
		dates[point.Date] = true
	}

	// Create time index.
	var sortedDates []string
	for date := range dates {
		sortedDates = append(sortedDates, date)
	}
	sort.StringSlice(sortedDates).Sort()

	var dateTimes []time.Time
	for _, date := range sortedDates {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}

		dateTimes = append(dateTimes, t)
	}

	var fields []*data.Field
	// Create fields for each series.
	for name, series := range series {
		values := make([]int64, len(sortedDates))
		for i, date := range sortedDates {
			if val, ok := series[date]; ok {
				values[i] = val
			}
		}
		fields = append(fields, data.NewField(name, nil, values))
	}

	// Sort the fields by name for more consistent styling in Grafana.
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return data.NewFrame(figure.Data.Type,
		// Prepend the time field before of the value fields.
		append([]*data.Field{
			data.NewField("time", nil, dateTimes),
		}, fields...)...,
	), nil
}
