import * as React from "react";
import { 
    Edit,
    SelectInput,
    TextInput,
    Create,
    SimpleForm,
    BooleanInput,
    NumberInput,
    DateTimeInput
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
        <TextInput disabled source="id" helperText="Job id. Must be unique, it's a copy of name." />
        <TextInput source="name" helperText="Job name. Must be unique, acts as the id." />
        <TextInput source="displayname" helperText="Display name of the job. If present, displayed instead of the name." />
        <TextInput source="timezone" helperText="The timezone where the cron expression will be evaluated in." />
        <TextInput source="schedule" helperText="Cron expression for the job. When to run the job." />
        <TextInput source="owner" helperText="Arbitrary string indicating the owner of the job." />
        <TextInput source="owner_email" helperText="Email address to use for notifications."/>
        <TextInput source="parent_job" helperText="Job id of job that this job is dependent upon." />
        <BooleanInput source="ephemeral" helperText="Delete the job after the first successful execution." />
        <DateTimeInput source="expires_at" helperText="The job will not be executed after this time." />
        <SelectInput source="concurrency" 
            choices={[
                { id: 'allow', name: 'Allow' },
                { id: 'forbid', name: 'Forbid' },
            ]}
            helperText="Concurrency policy for this job (allow, forbid)."
        />
        <JsonInput
            source="processors"
            reactJsonOptions={{
                name: null,
                collapsed: false,
                enableClipboard: true,
                displayDataTypes: false,
                src: {},
            }}
            helperText="Processor plugins to use for this job."
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
            helperText="Tags of the target servers to run this job against."
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
            helperText="Job metadata describes the job and allows filtering from the API."
        />
        <TextInput source="executor" helperText="Executor plugin to be used in this job." />
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
            helperText="Configuration arguments for the specific executor."
        />
        <BooleanInput source="disabled" helperText="Is this job disabled?" />
        <NumberInput source="retries" helperText="Number of times to retry a job that failed an execution." />
    </SimpleForm>
);
