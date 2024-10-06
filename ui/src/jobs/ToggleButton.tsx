import { useState } from 'react';
import { useNotify, useRefresh, Button, useRecordContext } from 'react-admin';
import { apiUrl, httpClient } from '../dataProvider';
import ToggleIcon from '@mui/icons-material/Pause';

const ToggleButton = () => {
    const record = useRecordContext();
    const refresh = useRefresh();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        httpClient(`${apiUrl}/jobs/${record.id}/toggle`, { method: 'POST' }) 
            .then(() => {
                notify('Job toggled');
                refresh();
            })
            .catch((e) => {
                notify('Error on toggle job', { type: 'warning' })
            })
            .finally(() => {
                setLoading(false);
            });
    };
    return (
        <Button 
            label="Toggle"
            onClick={handleClick}
            disabled={loading}
        >
            <ToggleIcon/>
        </Button>
    );
};

export default ToggleButton;
