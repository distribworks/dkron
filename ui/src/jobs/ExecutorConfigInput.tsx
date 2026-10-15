import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import FormHelperText from '@mui/material/FormHelperText';
import Typography from '@mui/material/Typography';
import { Form } from '@rjsf/mui';
import validator from '@rjsf/validator-ajv8';
import { RJSFSchema, UiSchema } from '@rjsf/utils';
import { InputHelperText, required, useInput } from 'react-admin';
import { useFormContext, useWatch } from 'react-hook-form';
import { apiUrl, httpClient } from '../dataProvider';
import { EDITOR_ERRORS_FIELD, StringMap, StringMapInput } from './StringMapInput';

type ExecutorPlugin = {
    id: string;
    name: string;
    schema?: RJSFSchema;
};

const isSecretKey = (key: string) => /(password|secret|token|private.?key|credential)/i.test(key);

export const schemaEditorMode = (executor: string | undefined, plugins: ExecutorPlugin[]) => {
    if (!executor) return 'json';
    const plugin = plugins.find(candidate => candidate.name === executor || candidate.id === executor);
    return plugin?.schema ? 'schema' : 'json';
};

export const clearExecutorConfigEditorError = (
    setValue: ReturnType<typeof useFormContext>['setValue'],
) => {
    setValue(`${EDITOR_ERRORS_FIELD}.executor_config`, undefined, {
        shouldDirty: false,
        shouldTouch: false,
        shouldValidate: false,
    });
};

const parseSchemaValue = (value: unknown, propertySchema: RJSFSchema): unknown => {
    if (typeof value !== 'string') {
        return value;
    }

    const type = Array.isArray(propertySchema.type) ? propertySchema.type[0] : propertySchema.type;
    if (type === 'boolean') {
        if (value === 'true') return true;
        if (value === 'false') return false;
    }
    if (type === 'number' || type === 'integer') {
        const parsed = Number(value);
        return Number.isNaN(parsed) ? value : parsed;
    }
    if (type === 'array' || type === 'object') {
        try {
            return JSON.parse(value);
        } catch {
            return value;
        }
    }
    return value;
};

export const toSchemaFormData = (value: StringMap | undefined, schema: RJSFSchema) => {
    const properties = (schema.properties || {}) as Record<string, RJSFSchema>;
    return Object.entries(value || {}).reduce<Record<string, unknown>>((acc, [key, currentValue]) => {
        acc[key] = parseSchemaValue(currentValue, properties[key] || {});
        return acc;
    }, {});
};

export const toExecutorConfig = (formData: Record<string, unknown> | undefined) =>
    Object.entries(formData || {}).reduce<StringMap>((acc, [key, value]) => {
        if (value === undefined || value === null) {
            return acc;
        }
        if (typeof value === 'string') {
            acc[key] = value;
            return acc;
        }
        if (typeof value === 'boolean' || typeof value === 'number') {
            acc[key] = String(value);
            return acc;
        }
        acc[key] = JSON.stringify(value);
        return acc;
    }, {});

export const validateExecutorConfigWithSchema = (value: StringMap | undefined, schema: RJSFSchema) => {
    const formData = toSchemaFormData(value, schema);
    const result = validator.validateFormData(formData, schema);
    if (!result.errors.length) return undefined;

    return result.errors
        .map(error => error.stack || [error.property, error.message].filter(Boolean).join(' '))
        .join('; ');
};

export const buildExecutorConfigUiSchema = (schema: RJSFSchema): UiSchema => {
    const properties = (schema.properties || {}) as Record<string, RJSFSchema>;
    return Object.keys(properties).reduce<UiSchema>((uiSchema, key) => {
        if (isSecretKey(key)) {
            uiSchema[key] = { 'ui:widget': 'password' };
        }
        return uiSchema;
    }, {
        'ui:submitButtonOptions': { norender: true },
        'ui:globalOptions': { label: true },
    });
};

const JsonExecutorConfigInput = () => (
    <StringMapInput
        source="executor_config"
        label="Executor configuration"
        helperText="Executor is not registered or has no schema. Configuration values are strings. Recognized secrets are masked in the form."
        emptyText="No executor configuration has been added."
        advancedLabel="Edit executor config as JSON"
    />
);

const SchemaExecutorConfigInput = ({ schema }: { schema: RJSFSchema }) => {
    const mounted = useRef(false);
    const { field, fieldState } = useInput({
        source: 'executor_config',
        validate: required(),
        defaultValue: {},
    });
    const { clearErrors, setError, setValue, trigger } = useFormContext();
    const formData = useMemo(
        () => toSchemaFormData(field.value, schema),
        [field.value, schema]
    );
    const uiSchema = useMemo(() => buildExecutorConfigUiSchema(schema), [schema]);

    useEffect(() => {
        mounted.current = true;
        return () => { mounted.current = false; };
    }, []);

    const applyValidation = (nextValue: StringMap) => {
        const message = validateExecutorConfigWithSchema(nextValue, schema);
        setValue(`${EDITOR_ERRORS_FIELD}.executor_config`, message, {
            shouldDirty: false,
            shouldTouch: false,
            shouldValidate: false,
        });
        if (message) {
            setError('executor_config', { type: 'validate', message });
        } else {
            clearErrors('executor_config');
        }
        void trigger('executor_config');
        return !message;
    };

    useEffect(() => {
        clearExecutorConfigEditorError(setValue);
        applyValidation((field.value ?? {}) as StringMap);
    }, [schema]);

    return (
        <Box sx={{ width: '100%' }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
                Executor configuration
            </Typography>
            <Form
                schema={schema}
                uiSchema={uiSchema}
                validator={validator}
                formData={formData}
                tagName="div"
                liveValidate
                noHtml5Validate
                showErrorList={false}
                onBlur={() => field.onBlur()}
                onChange={event => {
                    // RJSF can emit defaults from its constructor, before react-admin's
                    // useInput event handler is ready to be called.
                    if (!mounted.current) return;
                    const nextValue = toExecutorConfig(event.formData);
                    field.onChange(nextValue);
                    applyValidation(nextValue);
                }}
            >
                <></>
            </Form>
            <FormHelperText error={!!fieldState.error}>
                <InputHelperText
                    error={fieldState.error?.message}
                    helperText="Configuration arguments for the selected executor."
                />
            </FormHelperText>
        </Box>
    );
};

export const ExecutorConfigInput = () => {
    const executor = useWatch({ name: 'executor' });
    const [plugins, setPlugins] = useState<ExecutorPlugin[]>([]);
    const [metadataStatus, setMetadataStatus] = useState<'loading' | 'loaded' | 'error'>('loading');
    const previousExecutor = useRef(executor);
    const executorConfigs = useRef<Record<string, StringMap>>({});
    const { clearErrors, getValues, setValue, trigger } = useFormContext();

    const selectedPlugin = plugins.find(plugin => plugin.name === executor || plugin.id === executor);
    const schema = selectedPlugin?.schema;

    useLayoutEffect(() => {
        if (previousExecutor.current === executor) return;
        const previous = previousExecutor.current;
        const currentConfig = (getValues('executor_config') ?? {}) as StringMap;
        if (previous) {
            executorConfigs.current[previous] = { ...currentConfig };
        }
        previousExecutor.current = executor;

        // executor_config is the API's single map for the selected executor.
        // Keep separate in-memory drafts so fields from the previous executor
        // are neither displayed nor submitted for the new one.
        const nextConfig = executor ? executorConfigs.current[executor] ?? {} : {};
        setValue('executor_config', { ...nextConfig }, {
            shouldDirty: false,
            shouldTouch: false,
            shouldValidate: false,
        });

        // Clear the outgoing editor's errors before the new schema editor's
        // passive validation effect runs.
        clearExecutorConfigEditorError(setValue);
        clearErrors('executor_config');
        // Legacy executors have no schema editor to validate on mount; run the
        // form-level required-field checks against the newly selected config.
        if (!schema) void trigger('executor_config');
    }, [clearErrors, executor, getValues, schema, setValue, trigger]);

    const loadExecutorMetadata = useCallback(() => {
        setMetadataStatus('loading');
        httpClient(`${apiUrl}/plugins/executors`)
            .then(({ json }) => {
                setPlugins(json);
                setMetadataStatus('loaded');
            })
            .catch(() => setMetadataStatus('error'));
    }, []);

    useEffect(() => {
        loadExecutorMetadata();
    }, [loadExecutorMetadata]);

    useEffect(() => {
        if (metadataStatus === 'loaded') {
            clearErrors('executor_config');
            void trigger('executor_config');
            return;
        }
        setValue(`${EDITOR_ERRORS_FIELD}.executor_config`, undefined, {
            shouldDirty: false,
            shouldTouch: false,
            shouldValidate: false,
        });
        setError('executor_config', {
            type: 'validate',
            message: metadataStatus === 'loading'
                ? 'Loading executor configuration metadata.'
                : 'Could not load executor configuration metadata. Retry by reloading the page.',
        });
    }, [clearErrors, metadataStatus, setError, setValue, trigger]);

    if (metadataStatus !== 'loaded') {
        return (
            <Box>
                <FormHelperText error={metadataStatus === 'error'} role={metadataStatus === 'error' ? 'alert' : undefined}>
                    {metadataStatus === 'loading'
                        ? 'Loading executor configuration metadata…'
                        : 'Could not load executor configuration metadata.'}
                </FormHelperText>
                {metadataStatus === 'error' && (
                    <Button type="button" size="small" onClick={loadExecutorMetadata}>
                        Retry
                    </Button>
                )}
            </Box>
        );
    }

    if (!executor || !schema) {
        return <JsonExecutorConfigInput key={executor ?? 'no-executor'} />;
    }

    return <SchemaExecutorConfigInput key={executor} schema={schema} />;
};
