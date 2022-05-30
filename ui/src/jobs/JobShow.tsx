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
    useNotify, fetchStart, fetchEnd, Button,
} from 'react-admin';
import ToggleButton from "./ToggleButton"
import RunButton from "./RunButton"
import { JsonField } from "react-admin-json-view";
import ZeroDateField from "./ZeroDateField";
import JobIcon from '@material-ui/icons/Update';
import FullIcon from '@material-ui/icons/BatteryFull';
import { Tooltip } from '@material-ui/core';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { apiUrl } from '../dataProvider';

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

const FullButton = ({record}: any) => {
    const dispatch = useDispatch();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        dispatch(fetchStart()); // start the global loading indicator 
        fetch(`${apiUrl}/jobs/${record.job_name}/executions/${record.id}`)
            .then((response) => {
                if (response.ok) {
                    notify('Success loading full output');
                    return response.json()
                }
                throw response
            })
            .then((data) => {
                record.output_truncated = false
                record.output = data.output
            })
            .catch((e) => {
                notify('Error on loading full output', 'warning')
            })
            .finally(() => {
                setLoading(false);
                dispatch(fetchEnd()); // stop the global loading indicator
            });
    };
    return (
        <Button 
            label="Load full output"
            onClick={handleClick}
            disabled={loading}
        >
            <FullIcon/>
        </Button>
    );
};

const SpecialOutputPanel = ({ id, record, resource }: any) => {
    return (
        <div className="execution-output">
            {record.output_truncated ? <div><FullButton record={record} /></div> : ""}
            {record.output || "Nothing to show"}
        </div>
    );
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
                    }}
                />
            </Tab>
            <Tab label="executions" path="executions">
                <ReferenceManyField reference="executions" target="jobs" label="Executions">
                    <Datagrid rowClick="expand" isRowSelectable={ record => false } expand={<SpecialOutputPanel {...props} />}>
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
