import { Chip } from '@mui/material';

export const TagsField = ({ record }: any) => {
    if (record === undefined) {
		return null
    } else {
        return <ul>
            {Object.keys(record.Tags).map(key => (
                <Chip label={key+": "+record.Tags[key]} />
            ))}
        </ul>
    }
};
TagsField.defaultProps = {
    addLabel: true
};
