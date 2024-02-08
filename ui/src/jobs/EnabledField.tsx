import SuccessIcon from '@mui/icons-material/CheckCircle';
import FailedIcon from '@mui/icons-material/Cancel';
import { Tooltip } from '@mui/material';
import { useRecordContext } from 'react-admin';

const EnabledField = () => {
    const record = useRecordContext();

    if (record.disabled) {
        return <Tooltip title="Disabled"><FailedIcon htmlColor="red" /></Tooltip>
    } else {
        return <Tooltip title="Enabled"><SuccessIcon htmlColor="green" /></Tooltip>
    } 
};

export default EnabledField;
