import * as React from 'react';
import { useState } from 'react';
import { useNotify, useRefresh, Button } from 'react-admin';
import { apiUrl } from '../dataProvider';
import RunIcon from '@mui/icons-material/PlayArrow';

const RunButton = ({record}: any) => {
    const refresh = useRefresh();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        fetch(`${apiUrl}/jobs/${record.id}`, { method: 'POST' })
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
