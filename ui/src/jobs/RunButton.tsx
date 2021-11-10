import * as React from 'react';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNotify, useRefresh, fetchStart, fetchEnd, Button } from 'react-admin';
import { apiUrl } from '../dataProvider';
import RunIcon from '@material-ui/icons/PlayArrow';

const RunButton = ({record}: any) => {
    const dispatch = useDispatch();
    const refresh = useRefresh();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        dispatch(fetchStart()); // start the global loading indicator 
        fetch(`${apiUrl}/jobs/${record.id}`, { method: 'POST' })
            .then(() => {
                notify('Success running job');
                refresh();
            })
            .catch((e) => {
                notify('Error on running job', 'warning')
            })
            .finally(() => {
                setLoading(false);
                dispatch(fetchEnd()); // stop the global loading indicator
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
