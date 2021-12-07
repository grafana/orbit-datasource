import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { OrbitDataSourceOptions, OrbitQuery } from './types';

export class DataSource extends DataSourceWithBackend<OrbitQuery, OrbitDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<OrbitDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: OrbitQuery, scopedVars: ScopedVars): OrbitQuery {
    const templateSrv = getTemplateSrv();

    return {
      ...query,
      analyze: templateSrv.replace(query.analyze, scopedVars),
      affiliation: query.affiliation ? templateSrv.replace(query.affiliation, scopedVars) : query.affiliation,
    };
  }
}
