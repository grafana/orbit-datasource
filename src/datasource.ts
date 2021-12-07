import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { OrbitDataSourceOptions, OrbitQuery } from './types';

export class DataSource extends DataSourceWithBackend<OrbitQuery, OrbitDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<OrbitDataSourceOptions>) {
    super(instanceSettings);
  }
}
