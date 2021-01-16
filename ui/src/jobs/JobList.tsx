import * as React from "react";
import { 
    Datagrid,
    TextField,
    NumberField,
    DateField,
    BooleanField,
    EditButton,
    Filter,
    TextInput,
    List,
    SelectInput,
    BulkDeleteButton
} from 'react-admin';
import { Fragment } from 'react';
import BulkRunButton from "./BulkRunButton"
import BulkToggleButton from "./BulkToggleButton"

const JobFilter = (props: any) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <SelectInput source="status" choices={[
            { id: 'success', name: 'Success' },
            { id: 'failed', name: 'Failed' },
        ]} />
    </Filter>
);

const StatusField = (props: any) => {
    return <TextField {...props} />;
};

const JobBulkActionButtons = (props: any) => (
    <Fragment>
        <BulkRunButton {...props} />
        <BulkToggleButton {...props} />
        <BulkDeleteButton {...props} />
    </Fragment>
);

const JobList = (props: any) => (
    <List {...props} filters={<JobFilter />} bulkActionButtons={<JobBulkActionButtons />}>
        <Datagrid rowClick="show" style={{tableLayout: 'fixed'}}>
            <TextField source="id" />
            <TextField source="displayname" label="Display name" />
            <TextField source="schedule" />
            <NumberField source="success_count" />
            <NumberField source="error_count" />
            <DateField source="last_success" showTime />
            <DateField source="last_error" showTime />
            <BooleanField source="disabled" sortable={false} />
            <NumberField source="retries" sortable={false} />
            <StatusField source="status" sortable={false} />
            <DateField source="next" showTime />
            <EditButton/>
        </Datagrid>
    </List>
);

export default JobList;
