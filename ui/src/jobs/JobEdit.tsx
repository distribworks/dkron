import { ReactNode } from 'react';
import {
    AutocompleteInput,
    BooleanInput,
    Button as RaButton,
    Create,
    DateTimeInput,
    DeleteButton,
    Edit,
    ListButton,
    NumberInput,
    ReferenceInput,
    SaveButton,
    SelectInput,
    ShowButton,
    SimpleForm,
    TextInput,
    Toolbar,
    useInput,
    useRecordContext,
} from 'react-admin';
import { useFormContext, useFormState, useWatch } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import AddTaskIcon from '@mui/icons-material/AddTask';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import CodeIcon from '@mui/icons-material/Code';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import NotificationsOutlinedIcon from '@mui/icons-material/NotificationsOutlined';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import RouteIcon from '@mui/icons-material/Route';
import ScheduleIcon from '@mui/icons-material/Schedule';
import SettingsIcon from '@mui/icons-material/Settings';
import TuneIcon from '@mui/icons-material/Tune';
import {
    Autocomplete,
    Box,
    Card,
    CardContent,
    Chip,
    Divider,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { ProcessorsInput } from './ProcessorsInput';
import { EDITOR_ERRORS_FIELD, StringMap, StringMapInput } from './StringMapInput';
import { ExecutorConfigInput } from './ExecutorConfigInput';
import { pagePadding } from '../layout/Page';

const executorChoices = [
    'shell', 'http', 'grpc', 'kafka', 'nats', 'rabbitmq', 'gcppubsub',
];

const getTimezones = () => {
    const intl = Intl as typeof Intl & { supportedValuesOf?: (key: 'timeZone') => string[] };
    const zones = intl.supportedValuesOf?.('timeZone') ?? [];
    return Array.from(new Set(['UTC', ...zones]));
};

const timezoneChoices = getTimezones();

const SectionCard = ({
    icon: Icon,
    title,
    description,
    children,
}: {
    icon: typeof SettingsIcon;
    title: string;
    description: string;
    children: ReactNode;
}) => (
    <Card sx={{ mb: 3, width: '100%', boxSizing: 'border-box' }}>
        <CardContent>
            <Stack direction="row" spacing={1.5} alignItems="center">
                <Box
                    sx={{
                        width: 38,
                        height: 38,
                        borderRadius: 1.5,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: 'primary.main',
                        backgroundColor: 'rgba(49, 130, 206, 0.1)',
                    }}
                >
                    <Icon fontSize="small" />
                </Box>
                <Box>
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                        {title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {description}
                    </Typography>
                </Box>
            </Stack>
            <Divider sx={{ my: 2.5 }} />
            {children}
        </CardContent>
    </Card>
);

const FieldGrid = ({ children }: { children: ReactNode }) => (
    <Box
        sx={{
            display: 'grid',
            gridTemplateColumns: { xs: '1fr', md: 'repeat(2, minmax(0, 1fr))' },
            columnGap: 2.5,
            rowGap: 0.5,
            '& .MuiFormControl-root': { width: '100%' },
        }}
    >
        {children}
    </Box>
);

const JobFormHeader = ({ creating }: { creating: boolean }) => {
    const record = useRecordContext();
    const name = useWatch({ name: 'name' });
    const displayName = useWatch({ name: 'displayname' });
    const schedule = useWatch({ name: 'schedule' });
    const executor = useWatch({ name: 'executor' });

    return (
        <Box sx={{ mb: 3 }}>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
                <Box
                    sx={{
                        width: 52,
                        height: 52,
                        flex: '0 0 auto',
                        borderRadius: 2,
                        background: 'linear-gradient(135deg, #1a365d 0%, #2c5282 100%)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: 'white',
                    }}
                >
                    {creating ? <AddTaskIcon /> : <TuneIcon />}
                </Box>
                <Box sx={{ minWidth: 0 }}>
                    <Typography variant="h5" sx={{ fontWeight: 650 }}>
                        {creating ? 'Create job' : `Edit ${displayName || name || record?.name || 'job'}`}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {creating
                            ? 'Define when, where and how this job should run.'
                            : 'Update the schedule and execution settings without changing the job identity.'}
                    </Typography>
                </Box>
                <Stack
                    spacing={1}
                    alignItems={{ sm: 'flex-end' }}
                    sx={{ ml: { sm: 'auto' }, flex: '0 0 auto' }}
                >
                    <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap' }}>
                        {executor && <Chip size="small" icon={<CodeIcon />} label={executor} />}
                        {schedule && <Chip size="small" icon={<ScheduleIcon />} label={schedule} />}
                    </Stack>
                    <Stack direction="row" spacing={0.5}>
                        <ListButton />
                        {!creating && <ShowButton />}
                    </Stack>
                </Stack>
            </Stack>
        </Box>
    );
};

const ExecutorInput = () => {
    const { field, fieldState: { error } } = useInput({ source: 'executor' });
    return (
        <Autocomplete
            freeSolo
            options={executorChoices}
            value={field.value || null}
            inputValue={field.value || ''}
            onChange={(_, value) => field.onChange(value ?? '')}
            onInputChange={(_, value, reason) => {
                if (reason === 'input' || reason === 'clear') field.onChange(value);
            }}
            renderInput={params => (
                <TextField
                    {...params}
                    required
                    label="Executor"
                    error={Boolean(error)}
                    helperText={error?.message || 'Select a built-in executor or enter a custom plugin name.'}
                />
            )}
        />
    );
};

const TimezoneInput = () => {
    const { field, fieldState: { error } } = useInput({ source: 'timezone', defaultValue: '' });
    return (
        <Autocomplete
            options={timezoneChoices}
            value={field.value || null}
            onChange={(_, value) => field.onChange(value ?? '')}
            renderInput={params => (
                <TextField
                    {...params}
                    label="Timezone"
                    error={Boolean(error)}
                    helperText={error?.message || 'Leave empty to use the Dkron server timezone.'}
                />
            )}
        />
    );
};

const ScheduleInput = () => {
    const { setValue } = useFormContext();
    const parentJob = useWatch({ name: 'parent_job' });
    const presets = ['@manually', '@minutely', '@every 5m', '0 0 * * * *'];
    return (
        <Box sx={{ gridColumn: '1 / -1' }}>
            <TextInput
                source="schedule"
                label="Schedule"
                fullWidth
                helperText={parentJob
                    ? 'Optional when a parent job is selected. Add a schedule to allow both trigger paths.'
                    : 'Six-field cron expression including seconds, or a Dkron descriptor.'}
            />
            <Stack direction="row" spacing={1} sx={{ mt: -1, mb: 2, flexWrap: 'wrap', gap: 0.5 }}>
                <Typography variant="caption" color="text.secondary" sx={{ alignSelf: 'center', mr: 0.5 }}>
                    Presets:
                </Typography>
                {presets.map(preset => (
                    <Chip
                        key={preset}
                        size="small"
                        variant="outlined"
                        label={preset}
                        onClick={() => setValue('schedule', preset, { shouldDirty: true, shouldValidate: true })}
                    />
                ))}
            </Stack>
        </Box>
    );
};

const JobFormToolbar = ({ creating }: { creating: boolean }) => {
    const navigate = useNavigate();
    const { isDirty } = useFormState();
    const editorErrors = useWatch({ name: EDITOR_ERRORS_FIELD }) as Record<string, string | undefined> | undefined;
    const hasEditorErrors = Object.values(editorErrors ?? {}).some(Boolean);
    return (
        <Toolbar
            sx={{
                position: 'sticky',
                bottom: 0,
                zIndex: 5,
                borderTop: '1px solid',
                borderColor: 'divider',
                gap: 1,
                flexWrap: 'wrap',
            }}
        >
            <SaveButton
                label={creating ? 'Create job' : 'Save changes'}
                disabled={hasEditorErrors}
            />
            <RaButton label="Cancel" onClick={() => navigate(-1)}>
                <ArrowBackIcon />
            </RaButton>
            {isDirty && (
                <Chip
                    size="small"
                    color="warning"
                    variant="outlined"
                    label="Unsaved changes"
                    sx={{ ml: 'auto' }}
                />
            )}
            {!creating && <DeleteButton mutationMode="pessimistic" sx={{ ml: 'auto' }} />}
        </Toolbar>
    );
};

const isValidSlug = (value: string) => Array.from(value).every(character => {
    if (/[0-9_-]/.test(character)) return true;
    return character.toLowerCase() === character && character.toUpperCase() !== character;
});

const validateJob = (values: Record<string, any>) => {
    const errors: Record<string, string> = {};
    if (!values.name) errors.name = 'Job name is required.';
    else if (!isValidSlug(values.name)) errors.name = 'Use lowercase letters, numbers, underscores or dashes.';
    if (!values.executor) errors.executor = 'Executor is required.';
    if (!values.schedule && !values.parent_job) errors.schedule = 'A schedule or parent job is required.';
    if (values.parent_job && values.parent_job === values.name) errors.parent_job = 'A job cannot depend on itself.';
    if (values.owner_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.owner_email)) {
        errors.owner_email = 'Enter a valid email address.';
    }
    if (values.retries !== undefined && (!Number.isInteger(Number(values.retries)) || Number(values.retries) < 0)) {
        errors.retries = 'Retries must be a non-negative whole number.';
    }
    if (values.starts_at && values.expires_at && new Date(values.starts_at) >= new Date(values.expires_at)) {
        errors.expires_at = 'Expiration must be after the start time.';
    }

    const editorErrors = values[EDITOR_ERRORS_FIELD] as Record<string, string | undefined> | undefined;
    Object.entries(editorErrors ?? {}).forEach(([source, message]) => {
        if (message) errors[source] = message;
    });
    return errors;
};

const prepareJobForSave = (values: Record<string, any>) => {
    const job = { ...values };
    delete job[EDITOR_ERRORS_FIELD];
    return job;
};

const defaultValues = {
    concurrency: 'allow',
    disabled: false,
    ephemeral: false,
    executor: 'shell',
    executor_config: {},
    metadata: {},
    processors: {},
    retries: 0,
    tags: {},
    timezone: '',
};

const JobForm = ({ creating }: { creating: boolean }) => (
    <SimpleForm
        defaultValues={creating ? defaultValues : undefined}
        validate={validateJob}
        warnWhenUnsavedChanges
        toolbar={<JobFormToolbar creating={creating} />}
        sx={{
            width: '100%',
            maxWidth: 1120,
            mx: 'auto',
            '& > .MuiCardContent-root': { p: 0 },
        }}
    >
        <JobFormHeader creating={creating} />

        <SectionCard
            icon={InfoOutlinedIcon}
            title="Identity and ownership"
            description="How operators identify this job and who receives notifications."
        >
            <FieldGrid>
                {!creating && (
                    <TextInput
                        source="id"
                        label="Job ID"
                        readOnly
                        helperText="The ID is derived from the immutable job name."
                    />
                )}
                <TextInput
                    source="name"
                    label="Job name"
                    readOnly={!creating}
                    helperText={creating
                        ? 'Unique identifier using lowercase letters, numbers, underscores or dashes.'
                        : 'Job names cannot be changed safely after creation.'}
                />
                <TextInput
                    source="displayname"
                    label="Display name"
                    helperText="Optional human-friendly name shown throughout the UI."
                />
                <TextInput
                    source="owner"
                    label="Owner"
                    readOnly
                    helperText="Set from the authenticated API identity."
                />
                <TextInput
                    source="owner_email"
                    label="Notification email"
                    type="email"
                    helperText="Email address used by configured notification integrations."
                />
            </FieldGrid>
        </SectionCard>

        <SectionCard
            icon={ScheduleIcon}
            title="Schedule and dependencies"
            description="Choose the trigger, timezone and active time window."
        >
            <FieldGrid>
                <ScheduleInput />
                <TimezoneInput />
                <ReferenceInput
                    source="parent_job"
                    reference="jobs"
                    sort={{ field: 'name', order: 'ASC' }}
                    perPage={100}
                >
                    <AutocompleteInput
                        label="Parent job"
                        optionText="name"
                        helperText="Run after this job completes successfully."
                    />
                </ReferenceInput>
                <DateTimeInput
                    source="starts_at"
                    label="Starts at"
                    helperText="Do not execute before this date and time."
                />
                <DateTimeInput
                    source="expires_at"
                    label="Expires at"
                    helperText="Do not execute after this date and time."
                />
            </FieldGrid>
        </SectionCard>

        <SectionCard
            icon={PlayCircleOutlineIcon}
            title="Execution"
            description="Configure the plugin, failure handling and concurrency behavior."
        >
            <FieldGrid>
                <ExecutorInput />
                <SelectInput
                    source="concurrency"
                    choices={[
                        { id: 'allow', name: 'Allow overlapping executions' },
                        { id: 'forbid', name: 'Forbid overlapping executions' },
                    ]}
                    helperText="Choose whether a new run may start while the previous run is active."
                />
                <NumberInput
                    source="retries"
                    min={0}
                    step={1}
                    helperText="Number of additional attempts after a failed execution."
                />
            </FieldGrid>
            <Box sx={{ mt: 2 }}>
                <ExecutorConfigInput />
            </Box>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={{ xs: 0, sm: 4 }} sx={{ mt: 2 }}>
                <BooleanInput
                    source="disabled"
                    label="Disabled"
                    helperText="Keep the job configured without scheduling new runs."
                />
                <BooleanInput
                    source="ephemeral"
                    label="Delete after success"
                    helperText="Remove this job after its first successful execution."
                />
            </Stack>
        </SectionCard>

        <SectionCard
            icon={RouteIcon}
            title="Targeting and output"
            description="Select target agents, attach metadata and process execution output."
        >
            <Stack spacing={3}>
                <StringMapInput
                    source="tags"
                    label="Target tags"
                    helperText="Only agents matching these tag values are eligible to run the job."
                    emptyText="No tag filter. Any eligible agent can run this job."
                />
                <Divider />
                <StringMapInput
                    source="metadata"
                    label="Metadata"
                    helperText="Searchable labels for API clients and operators."
                    emptyText="No metadata attached."
                />
                <Divider />
                <ProcessorsInput />
            </Stack>
        </SectionCard>

        <Stack direction="row" spacing={1} alignItems="center" sx={{ color: 'text.secondary', mb: 2 }}>
            <NotificationsOutlinedIcon fontSize="small" />
            <Typography variant="caption">
                Saving updates the job definition immediately on the Dkron leader.
            </Typography>
        </Stack>
    </SimpleForm>
);

export const JobEdit = () => (
    <Edit
        actions={false}
        mutationMode="pessimistic"
        transform={prepareJobForSave}
        sx={{
            '& .RaEdit-main': { p: pagePadding },
            '& .RaEdit-card': { width: '100%', maxWidth: 1184, mx: 'auto' },
        }}
    >
        <JobForm creating={false} />
    </Edit>
);

export const JobCreate = () => (
    <Create
        actions={false}
        mutationMode="pessimistic"
        transform={prepareJobForSave}
        sx={{
            '& .RaCreate-main': { p: pagePadding },
            '& .RaCreate-card': { width: '100%', maxWidth: 1184, mx: 'auto' },
        }}
    >
        <JobForm creating />
    </Create>
);
