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

const BulkRunButton = ({...props}: any) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const unselectAll = useUnselectAll;
    const [loading, setLoading] = useState(false);
    const { selectedIds } = useListContext();
    const runMany = () => {
        for(let id of selectedIds) {
            setLoading(true);
            httpClient(`${apiUrl}/jobs/${id}`, { method: 'POST' })
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
