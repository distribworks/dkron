import * as React from 'react';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNotify, useRedirect, fetchStart, fetchEnd, Button } from 'react-admin';
import { apiUrl } from '../dataProvider';
import ToggleIcon from '@material-ui/icons/Pause';

const ToggleButton = ({record}: any) => {
    const dispatch = useDispatch();
    const redirect = useRedirect();
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const handleClick = () => {
        setLoading(true);
        dispatch(fetchStart()); // start the global loading indicator 
        fetch(`${apiUrl}/jobs/${record.id}/toggle`, { method: 'POST' })
            .then(() => {
                notify('Job toggled');
                redirect('/jobs');
            })
            .catch((e) => {
                notify('Error on toggle job', 'warning')
            })
            .finally(() => {
                setLoading(false);
                dispatch(fetchEnd()); // stop the global loading indicator
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
