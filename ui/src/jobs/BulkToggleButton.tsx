import { useState } from 'react';
import {
    useNotify,
    Button,
    useUnselectAll,
    useRefresh,
} from 'react-admin';
import { apiUrl } from '../dataProvider';
import RunIcon from '@mui/icons-material/PlayArrow';

const BulkToggleButton = ({selectedIds}: any) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll;
    const [loading, setLoading] = useState(false);
    const toggleMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            fetch(`${apiUrl}/jobs/${id}/toggle`, { method: 'POST' })
                .then(() => {
                    notify('Job toggled');
                })
                .catch((e) => {
                    notify('Error on job toggle', { type: 'warning' })
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
