import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { OrbitDataSourceOptions, OrbitSecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<OrbitDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onWorkspacePathChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      workspacePath: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onAPITokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiToken: event.target.value,
      },
    });
  };

  onResetAPIToken = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiToken: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as OrbitSecureJsonData;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Workspace path"
            labelWidth={8}
            inputWidth={20}
            onChange={this.onWorkspacePathChange}
            value={jsonData.workspacePath || ''}
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.apiToken) as boolean}
              value={secureJsonData.apiToken || ''}
              label="API token"
              labelWidth={8}
              inputWidth={20}
              onReset={this.onResetAPIToken}
              onChange={this.onAPITokenChange}
            />
          </div>
        </div>
      </div>
    );
  }
}
