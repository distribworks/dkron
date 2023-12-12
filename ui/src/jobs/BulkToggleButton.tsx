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

const BulkToggleButton = ({selectedIds}: any) => {
    const dispatch = useDispatch();
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll();
    const [loading, setLoading] = useState(false);
    const toggleMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            dispatch(fetchStart()); // start the global loading indicator
            fetch(`${apiUrl}/jobs/${id}/toggle`, { method: 'POST' })
                .then(() => {
                    notify('Job toggled');
                })
                .catch((e) => {
                    notify('Error on job toggle', 'warning')
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
            label="Toggle"
            title='Enable/disable execution of selected jobs'
            onClick={toggleMany}
            disabled={loading}
        >
            <RunIcon/>
        </Button>
    );
};

export default BulkToggleButton;
