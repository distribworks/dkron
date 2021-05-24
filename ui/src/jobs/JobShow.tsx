import * as React from "react";
import { 
    Datagrid,
    TextField,
    NumberField,
    DateField,
    EditButton,
    BooleanField,
    TopToolbar,
    Show,
    TabbedShowLayout,
    Tab,
    ReferenceManyField,
} from 'react-admin';
import { OutputPanel } from "../executions/BusyList";
import ToggleButton from "./ToggleButton"
import RunButton from "./RunButton"
import { JsonField } from "react-admin-json-view";
import ZeroDateField from "./ZeroDateField";
import JobIcon from '@material-ui/icons/Update';
import { Tooltip } from '@material-ui/core';

const JobShowActions = ({ basePath, data, resource }: any) => (
    <TopToolbar>
        <RunButton record={data} />
        <ToggleButton record={data} />
        <EditButton basePath={basePath} record={data} />
    </TopToolbar>
);

const SuccessField = (props: any) => {
    return (props.record["finished_at"] === null ? <Tooltip title="Running"><JobIcon /></Tooltip> : <BooleanField {...props} />);
};

const JobShow = (props: any) => (
    <Show actions={<JobShowActions {...props}/>} {...props}>
        <TabbedShowLayout>
            <Tab label="summary">
                <TextField source="name" />
                <TextField source="timezone" />
                <TextField source="schedule" />
                <DateField label="Last success" source="last_success" showTime />
                <DateField source="last_error" showTime />
                <TextField source="status" />
                <TextField source="concurrency" />
                <BooleanField source="ephemeral" />
                <DateField source="expires_at" showTime />
                <DateField source="next" sortable={false} showTime />
                <JsonField
                    source="processors"
                    reactJsonOptions={{
                        // Props passed to react-json-view
                        name: null,
                        collapsed: false,
                        enableClipboard: false,
                        displayDataTypes: false,
                        src: {},
                    }}
                />
                <JsonField
                    source="tags"
                    reactJsonOptions={{
                        // Props passed to react-json-view
                        name: null,
                        collapsed: false,
                        enableClipboard: false,
                        displayDataTypes: false,
                        src: {},
                    }}
                />
                <JsonField
                    source="executor_config"
                    reactJsonOptions={{
                        // Props passed to react-json-view
                        name: null,
                        collapsed: false,
                        enableClipboard: false,
                        displayDataTypes: false,
                        src: {},
                    }}
                />
            </Tab>
            <Tab label="executions" path="executions">
                <ReferenceManyField reference="executions" target="jobs" label="Executions">
                    <Datagrid rowClick="expand" isRowSelectable={ record => false } expand={<OutputPanel {...props} />}>
                        <TextField source="id" />
                        <TextField source="group" sortable={false} />
                        <TextField source="job_name" sortable={false} />
                        <DateField source="started_at" showTime />
                        <ZeroDateField source="finished_at" showTime />
                        <TextField source="node_name" sortable={false} />
                        <SuccessField source="success" sortable={false} />
                        <NumberField source="attempt" />
                    </Datagrid>
                </ReferenceManyField>
            </Tab>
        </TabbedShowLayout>
    </Show>
);
export default JobShow;
