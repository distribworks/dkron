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
    SelectInput
} from 'react-admin';
import ToggleButton from "./ToggleButton"
import RunButton from "./RunButton"

const PostFilter = (props: any) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <SelectInput source="status" choices={[
            { id: 'success', name: 'Success' },
            { id: 'failed', name: 'Failed' },
        ]} />
    </Filter>
);

const JobList = (props: any) => (
    <List filters={<PostFilter />} {...props}>
        <Datagrid rowClick="show">
            <TextField source="id" />
            <TextField source="displayname" />
            <TextField source="schedule" />
            <NumberField source="success_count" />
            <NumberField source="error_count" />
            <DateField source="last_success" showTime />
            <DateField source="last_error" showTime />
            <BooleanField source="disabled" sortable={false} />
            <NumberField source="retries" sortable={false} />
            <TextField source="status" sortable={false} />
            <DateField source="next" showTime />
            <RunButton />
            <ToggleButton />
            <EditButton />
        </Datagrid>
    </List>
);

export default JobList;
