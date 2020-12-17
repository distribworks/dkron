import * as React from "react";
import { 
    Edit,
    SelectInput,
    TextInput,
    Create,
    SimpleForm,
    BooleanInput,
    NumberInput
} from 'react-admin';
import { JsonInput } from "react-admin-json-view";

export const JobEdit = (props: any) => (
    <Edit {...props}>
        <EditForm />
    </Edit>
);

export const JobCreate = (props: any) => (
    <Create {...props}>
        <EditForm />
    </Create>
);

const EditForm = (props: any) => (
    <SimpleForm  {...props}>
        <TextInput disabled source="id" />
        <TextInput source="name" />
        <TextInput source="schedule" />
        <TextInput source="displayname" />
        <TextInput source="owner" />
        <TextInput source="owner_email" />
        <TextInput source="parent_job" />
        <SelectInput source="concurrency" choices={[
            { id: 'allow', name: 'Allow' },
            { id: 'forbid', name: 'Forbid' },
        ]} />
        <JsonInput
            source="processors"
            // validate={required(){ return true }}
            reactJsonOptions={{
                // Props passed to react-json-view
                name: null,
                collapsed: false,
                enableClipboard: false,
                displayDataTypes: false,
                src: {},
            }}
        />
        <JsonInput
            source="tags"
            // validate={required(){ return true }}
            reactJsonOptions={{
                // Props passed to react-json-view
                name: null,
                collapsed: false,
                enableClipboard: false,
                displayDataTypes: false,
                src: {},
            }}
        />
        <JsonInput
            source="metadata"
            // validate={required(){ return true }}
            reactJsonOptions={{
                // Props passed to react-json-view
                name: null,
                collapsed: false,
                enableClipboard: false,
                displayDataTypes: false,
                src: {},
            }}
        />
        <TextInput source="executor" />
        <JsonInput
            source="executor_config"
            // validate={required(){ return true }}
            reactJsonOptions={{
                // Props passed to react-json-view
                name: null,
                collapsed: true,
                enableClipboard: false,
                displayDataTypes: false,
                src: {},
            }}
        />
        <BooleanInput source="disabled" />
        <NumberInput source="retries" />
    </SimpleForm>
);
