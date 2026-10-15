import { describe, expect, it } from 'vitest';
import {
    buildExecutorConfigUiSchema,
    clearExecutorConfigEditorError,
    schemaEditorMode,
    toExecutorConfig,
    toSchemaFormData,
    validateExecutorConfigWithSchema,
} from './ExecutorConfigInput';

describe('executor config schema helpers', () => {
    const schema = {
        type: 'object',
        required: ['url', 'retries'],
        properties: {
            url: { type: 'string', minLength: 1 },
            retries: { type: 'integer', minimum: 1 },
            password: { type: 'string' },
            saslPassword: { type: 'string' },
        },
    };

    it('validates custom executor requirements and constraints', () => {
        expect(validateExecutorConfigWithSchema({}, schema)).toContain("must have required property 'url'");
        expect(validateExecutorConfigWithSchema({ url: 'https://example.com', retries: '0' }, schema)).toContain('must be >= 1');
        expect(validateExecutorConfigWithSchema({ url: 'https://example.com', retries: '3' }, schema)).toBeUndefined();
    });

    it('converts existing string-map values to schema form data and back', () => {
        const formData = toSchemaFormData({ url: 'https://example.com', retries: '2' }, schema);

        expect(formData).toEqual({ url: 'https://example.com', retries: 2 });
        expect(toExecutorConfig(formData)).toEqual({ url: 'https://example.com', retries: '2' });
    });

    it('uses password widgets for secret-like schema keys', () => {
        const uiSchema = buildExecutorConfigUiSchema(schema);

        expect(uiSchema.password).toEqual({ 'ui:widget': 'password' });
        expect(uiSchema.saslPassword).toEqual({ 'ui:widget': 'password' });
        expect(uiSchema.url).toBeUndefined();
    });

    it('selects JSON fallback for unknown executors and schema mode for registered schemas', () => {
        const plugins = [
            { id: 'custom', name: 'custom', schema },
            { id: 'legacy', name: 'legacy' },
        ];

        expect(schemaEditorMode('custom', plugins)).toBe('schema');
        expect(schemaEditorMode('legacy', plugins)).toBe('json');
        expect(schemaEditorMode('missing', plugins)).toBe('json');
    });

    it('clears stale executor JSON editor errors on mode changes', () => {
        const calls: Array<[string, unknown]> = [];
        const setValue = (name: string, value: unknown) => calls.push([name, value]);

        clearExecutorConfigEditorError(setValue as any);

        expect(calls).toEqual([['__editor_errors.executor_config', undefined]]);
    });
});
