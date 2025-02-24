import { useState } from 'react';
import {
    useNotify,
    Button,
    useUnselectAll,
    useRefresh,
    useListContext,
} from 'react-admin';
import { apiUrl, httpClient } from '../dataProvider';
import RunIcon from '@mui/icons-material/PlayArrow';

const BulkToggleButton = ({...props}: any) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll;
    const { selectedIds } = useListContext();
    const [loading, setLoading] = useState(false);
    const toggleMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            httpClient(`${apiUrl}/jobs/${id}/toggle`, { method: 'POST' })
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
