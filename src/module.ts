import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { OrbitQuery, OrbitDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, OrbitQuery, OrbitDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
