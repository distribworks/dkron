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
        <TextInput source="displayname" />
        <TextInput source="timezone" />
        <TextInput source="schedule" />
        <TextInput source="owner" />
        <TextInput source="owner_email" />
        <TextInput source="parent_job" />
        <SelectInput source="concurrency" choices={[
            { id: 'allow', name: 'Allow' },
            { id: 'forbid', name: 'Forbid' },
        ]} />
        <JsonInput
            source="processors"
            reactJsonOptions={{
                name: null,
                collapsed: false,
                enableClipboard: true,
                displayDataTypes: false,
                src: {},
            }}
        />
        <JsonInput
            source="tags"
            reactJsonOptions={{
                name: null,
                collapsed: false,
                enableClipboard: true,
                displayDataTypes: false,
                src: {},
            }}
        />
        <JsonInput
            source="metadata"
            reactJsonOptions={{
                name: null,
                collapsed: false,
                enableClipboard: true,
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
