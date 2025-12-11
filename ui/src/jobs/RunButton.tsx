import { useState } from 'react';
import { useNotify, useRefresh, Button, useRecordContext } from 'react-admin';
import { apiUrl, httpClient } from '../dataProvider';
import RunIcon from '@mui/icons-material/PlayArrow';

const RunButton = () => {
    const record = useRecordContext();
    const refresh = useRefresh();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        if (!record) return;
        setLoading(true);
        httpClient(`${apiUrl}/jobs/${record.id}`, { method: 'POST' })
            .then(() => {
                notify('Success running job');
                refresh();
            })
            .catch((e) => {
                notify('Error on running job', { type: 'warning' })
            })
            .finally(() => {
                setLoading(false);
            });
    };
    return (
        <Button 
            label="Run"
            onClick={handleClick}
            disabled={loading}
        >
            <RunIcon/>
        </Button>
    );
};

export default RunButton;
