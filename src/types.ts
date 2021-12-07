import { DataQuery, DataSourceJsonData } from '@grafana/data';

export type Orbits = '1' | '2' | '3' | '4' | 'n';

export interface OrbitQuery extends DataQuery {
  analyze: string;
  interval?: string;
  affiliation?: string;
  newReturning?: string;
  groupBy?: string;
  groupLimit?: string;
  orbits?: Orbits[];
  activityTypes?: string[];
  cumulative?: boolean;
}

export const defaultQuery: Partial<OrbitQuery> = {
  analyze: 'members',
  interval: 'week',
  groupLimit: '10',
};

export interface OrbitDataSourceOptions extends DataSourceJsonData {
  workspacePath?: string;
}

export interface OrbitSecureJsonData {
  apiToken?: string;
}
