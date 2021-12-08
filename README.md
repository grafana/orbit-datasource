# Orbit data source for Grafana

[![CI](https://github.com/grafana/orbit-datasource/actions/workflows/ci.yml/badge.svg)](https://github.com/grafana/orbit-datasource/actions/workflows/ci.yml)
[![Release](https://github.com/grafana/orbit-datasource/actions/workflows/release.yml/badge.svg)](https://github.com/grafana/orbit-datasource/actions/workflows/release.yml)

A data source for [Orbit](https://orbit.love) that lets you monitor your community with Grafana.

## Prerequisites

- An [Orbit](https://orbit.love) account.

## Configure the data source

To configure the Orbit data source you need the _workspace path_ and a _API token_ from Orbit.

| Field          | Description                                       |
|----------------|---------------------------------------------------|
| Workspace path | Path to the Orbit workspace you want to monitor.  |
| API token      | Token used to authenticate against the Orbit API. |

You can get both of these fields by logging into your Orbit workspace.

- To get the workspace path, browse to **Settings** -> **Workspace Name** and copy the text in the **WORKSPACE PATH** text box.
- To get the API token, browse to **Account Settings** -> **API Token** and click **Copy**.
