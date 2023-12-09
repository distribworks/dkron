import * as React from 'react'
import { Chip } from '@mui/material';

export const TagsField = ({ record }: any) => (
    <ul>
        {Object.keys(record.Tags).map(key => (
            <Chip label={key+": "+record.Tags[key]} />
        ))}
    </ul>
)
TagsField.defaultProps = {
    addLabel: true
};
