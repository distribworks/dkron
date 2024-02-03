import * as React from 'react';
import { useState } from 'react';
import {
    useNotify,
    Button,
    useUnselectAll,
    useRefresh,
} from 'react-admin';
import { apiUrl } from '../dataProvider';
import RunIcon from '@mui/icons-material/PlayArrow';

const BulkRunButton = ({selectedIds}: any) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll;
    const [loading, setLoading] = useState(false);
    const runMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            fetch(`${apiUrl}/jobs/${id}`, { method: 'POST' })
                .then(() => {
                    notify('Success running job');
                })
                .catch((e) => {
                    notify('Error on running job', { type: 'warning' })
                })
                .finally(() => {
                    setLoading(false);
                    refresh();
                });
        }
        unselectAll('jobs');
    };
    return (
        <Button 
            label="Run"
            onClick={runMany}
            disabled={loading}
        >
            <RunIcon/>
        </Button>
    );
};

export default BulkRunButton;
