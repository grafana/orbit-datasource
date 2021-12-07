import React from 'react';
import {
  InlineFieldRow,
  InlineField,
  Select,
  RadioButtonGroup,
  FilterPill,
  HorizontalGroup,
  AsyncMultiSelect,
} from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, OrbitDataSourceOptions, OrbitQuery, Orbits } from './types';
import { defaults } from 'lodash';
import { getBackendSrv } from '@grafana/runtime';
import { map } from 'rxjs/operators';

type Props = QueryEditorProps<DataSource, OrbitQuery, OrbitDataSourceOptions>;

export const QueryEditor = (props: Props) => {
  const { onChange, onRunQuery } = props;

  const query = defaults(props.query, defaultQuery);

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Analyze" labelWidth={14}>
          <RadioButtonGroup
            value={query.analyze}
            onChange={(value) => {
              onChange({ ...query, analyze: value ?? 'members' });
              onRunQuery();
            }}
            options={[
              { label: 'Members', value: 'members' },
              { label: 'Activities', value: 'activities' },
            ]}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Group by" labelWidth={14}>
          <Select
            width={18}
            value={query.groupBy}
            onChange={(value) => {
              onChange({ ...query, groupBy: value?.value });
              onRunQuery();
            }}
            isClearable={true}
            options={[
              { label: 'Activity type', value: 'activity_type' },
              { label: 'Affiliation', value: 'affiliation' },
              { label: 'Channel', value: 'channel' },
              { label: 'Company', value: 'company' },
              { label: 'Group', value: 'group' },
              { label: 'Location', value: 'location' },
              { label: 'Member', value: 'member' },
              { label: 'Member tag', value: 'member_tag' },
              { label: 'Orbit level', value: 'orbit_level' },
              { label: 'Title', value: 'title' },
            ]}
          />
        </InlineField>
        <InlineField label="Group limit" labelWidth={14}>
          <Select
            width={14}
            value={query.groupLimit}
            isClearable={true}
            onChange={(value) => {
              onChange({ ...query, groupLimit: value?.value });
              onRunQuery();
            }}
            options={[
              { label: '5', value: '5' },
              { label: '10', value: '10' },
              { label: '25', value: '25' },
              { label: '50', value: '50' },
            ]}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Interval" labelWidth={14}>
          <RadioButtonGroup
            value={query.interval ?? 'century'}
            onChange={(interval) => {
              onChange({ ...query, interval });
              onRunQuery();
            }}
            options={[
              { label: 'Default', value: 'century' },
              { label: 'Daily', value: 'day' },
              { label: 'Weekly', value: 'week' },
              { label: 'Monthly', value: 'month' },
              { label: 'Yearly', value: 'year' },
            ]}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Orbit levels" labelWidth={14}>
          <HorizontalGroup>
            <FilterPill
              label="1"
              selected={!!query.orbits?.includes('1')}
              onClick={() => {
                onChange({
                  ...query,
                  orbits: toggleOrbit(query.orbits ?? [], '1'),
                });
                onRunQuery();
              }}
            />
            <FilterPill
              label="2"
              selected={!!query.orbits?.includes('2')}
              onClick={() => {
                onChange({
                  ...query,
                  orbits: toggleOrbit(query.orbits ?? [], '2'),
                });
                onRunQuery();
              }}
            />
            <FilterPill
              label="3"
              selected={!!query.orbits?.includes('3')}
              onClick={() => {
                onChange({
                  ...query,
                  orbits: toggleOrbit(query.orbits ?? [], '3'),
                });
                onRunQuery();
              }}
            />
            <FilterPill
              label="4"
              selected={!!query.orbits?.includes('4')}
              onClick={() => {
                onChange({
                  ...query,
                  orbits: toggleOrbit(query.orbits ?? [], '4'),
                });
                onRunQuery();
              }}
            />
            <FilterPill
              label="None"
              selected={!!query.orbits?.includes('n')}
              onClick={() => {
                onChange({
                  ...query,
                  orbits: toggleOrbit(query.orbits ?? [], 'n'),
                });
                onRunQuery();
              }}
            />
          </HorizontalGroup>
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Affiliation" labelWidth={14}>
          <RadioButtonGroup
            value={query.affiliation ?? ''}
            onChange={(affiliation) => {
              onChange({ ...query, affiliation });
              onRunQuery();
            }}
            options={[
              { label: 'All', value: '' },
              { label: 'Members', value: 'member' },
              { label: 'Teammates', value: 'team' },
            ]}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="New/returning" labelWidth={14}>
          <RadioButtonGroup
            value={query.newReturning ?? ''}
            onChange={(newReturning) => {
              onChange({ ...query, newReturning });
              onRunQuery();
            }}
            options={[
              {
                label: 'All',
                value: '',
              },
              {
                label: 'New members',
                value: 'new',
                description: 'First activity matching the filters was in the selected timeframe.',
              },
              {
                label: 'Returning members',
                value: 'returning',
                description:
                  'Had a matching activity before the selected timeframe and at least one during the selected timeframe.',
              },
            ]}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Activity types" labelWidth={14}>
          <AsyncMultiSelect
            width={48}
            value={
              props.query.activityTypes?.map<SelectableValue<string>>((activityType) => ({
                label: activityType,
                value: activityType,
              })) ?? []
            }
            onChange={(values: SelectableValue<string>[]) => {
              onChange({ ...query, activityTypes: values.map((value) => value.value!) });
              onRunQuery();
            }}
            defaultOptions
            loadOptions={async (query) => {
              const options = await getBackendSrv()
                .fetch<any>({
                  url: `/api/datasources/${props.datasource.id}/resources`,
                })
                .pipe(
                  map((response) => {
                    return response.data.map((activityType: any) => ({
                      label: activityType.attributes.short_name,
                      value: activityType.attributes.key,
                    }));
                  })
                )
                .toPromise();

              return options ?? [];
            }}
          />
        </InlineField>
      </InlineFieldRow>
    </>
  );
};

function toggleOrbit(orbits: Orbits[], orbit: Orbits): Orbits[] {
  const res = new Set<Orbits>(orbits);

  if (res.has(orbit)) {
    res.delete(orbit);
  } else {
    res.add(orbit);
  }

  return Array.from(res).sort();
}
