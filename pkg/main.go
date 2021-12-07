package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/orbit-datasource/pkg/plugin"
)

const pluginID = "marcusolsson-orbit-datasource"

func main() {
	backend.SetupPluginEnvironment(pluginID)

	if err := datasource.Manage(pluginID, plugin.NewOrbitDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
