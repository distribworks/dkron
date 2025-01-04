import { Chip } from '@mui/material';
import { useRecordContext } from 'react-admin';

export const TagsField = () => {
    const record = useRecordContext();
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
