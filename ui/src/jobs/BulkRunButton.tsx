import * as React from 'react';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import {
    useNotify,
    fetchStart,
    fetchEnd,
    Button,
    useUnselectAll,
    useRefresh,
} from 'react-admin';
import { apiUrl } from '../dataProvider';
import RunIcon from '@material-ui/icons/PlayArrow';

const BulkRunButton = ({selectedIds}: any) => {
    const dispatch = useDispatch();
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll();
    const [loading, setLoading] = useState(false);
    const runMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            dispatch(fetchStart()); // start the global loading indicator
            fetch(`${apiUrl}/jobs/${id}`, { method: 'POST' })
                .then(() => {
                    notify('Success running job');
                })
                .catch((e) => {
                    notify('Error on running job', 'warning')
                })
                .finally(() => {
                    setLoading(false);
                    refresh();
                    dispatch(fetchEnd()); // stop the global loading indicator
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
