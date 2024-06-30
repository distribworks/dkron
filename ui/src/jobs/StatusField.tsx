import SuccessIcon from '@mui/icons-material/CheckCircle';
import FailedIcon from '@mui/icons-material/Cancel';
import UntriggeredIcon from '@mui/icons-material/Timer';
import { Tooltip } from '@mui/material';
import { useRecordContext } from 'react-admin';

const StatusField = () => {
	const record = useRecordContext();
	
	if (record.status === 'success') {
		return <Tooltip title="Success"><SuccessIcon htmlColor="green" /></Tooltip> 
	} else if (record.status === 'failed') {
		return <Tooltip title="Error"><FailedIcon htmlColor="red" /></Tooltip>
	} else {
		return <Tooltip title="Waiting to Run"><UntriggeredIcon htmlColor="blue" /></Tooltip>
	}
};

export default StatusField;
