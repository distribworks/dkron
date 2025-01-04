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
    useNotify,
    Button,
    useRecordContext
} from 'react-admin';
import ToggleButton from "./ToggleButton"
import RunButton from "./RunButton"
import { JsonField } from "react-admin-json-view";
import ZeroDateField from "./ZeroDateField";
import JobIcon from '@mui/icons-material/Update';
import FullIcon from '@mui/icons-material/BatteryFull';
import { Tooltip } from '@mui/material';
import { useState } from 'react';
import { apiUrl, httpClient } from '../dataProvider';

// basePath={basePath}
const JobShowActions = ({ basePath, data, resource }: any) => (
    <TopToolbar>
        <RunButton />
        <ToggleButton />
        <EditButton record={data} />
    </TopToolbar>
);

const SuccessField = (props: any) => {
    return (
        props.record !== undefined && props.record["finished_at"] === null ? 
            <Tooltip title="Running"><JobIcon /></Tooltip> : <BooleanField {...props} />
    );
};

const FullButton = ({record}: any) => {
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        httpClient(`${apiUrl}/jobs/${record.job_name}/executions/${record.id}`)
            .then((response) => {
                if (response.status === 200) {
                    notify('Success loading full output');
                    return response.json
                }
                throw response
            })
            .then((data) => {
                record.output_truncated = false;
                record.output = data.output;
            })
            .catch((e) => {
                notify('Error on loading full output', { type: 'warning' })
            })
            .finally(() => {
                setLoading(false);
            });
    };

    if (record.output_truncated === false) return record.output;

    return (
        <Button 
            label="Load full output"
            onClick={handleClick}
            disabled={loading}
        >
            <FullIcon/>
        </Button>
    )
};

const SpecialOutputPanel = () => {
    const record = useRecordContext();
    if (!record) return null;
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
                        name: null,
                        collapsed: false,
                        enableClipboard: false,
                        displayDataTypes: false,
                    }}
                />
                <JsonField
                    source="tags"
                    reactJsonOptions={{
                        name: null,
                        collapsed: false,
                        enableClipboard: false,
                        displayDataTypes: false,
                    }}
                />
                <JsonField
                    source="metadata"
                    reactJsonOptions={{
                        name: null,
                        collapsed: false,
                        enableClipboard: true,
                        displayDataTypes: false,
                    }}
                />
                <TextField source="executor" />
                <JsonField
                    source="executor_config"
                    reactJsonOptions={{
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
