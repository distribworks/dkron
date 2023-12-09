import * as React from "react";
import SuccessIcon from '@mui/icons-material/CheckCircle';
import FailedIcon from '@mui/icons-material/Cancel';
import UntriggeredIcon from '@mui/icons-material/Timer';
import { Tooltip } from '@mui/material';

const StatusField = (props: any) => {
	if (props.record === undefined) {
		return null
	} else {
		if (props.record[props.source] === 'success') {
			return <Tooltip title="Success"><SuccessIcon htmlColor="green" /></Tooltip> 
		} else if (props.record[props.source] === 'failed') {
			return <Tooltip title="Error"><FailedIcon htmlColor="red" /></Tooltip>
		} else {
			return <Tooltip title="Waiting to Run"><UntriggeredIcon htmlColor="blue" /></Tooltip>
		}
	}
};

export default StatusField;
