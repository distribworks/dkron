import * as React from "react";
import SuccessIcon from '@material-ui/icons/CheckCircle';
import FailedIcon from '@material-ui/icons/Cancel';
import UntriggeredIcon from '@material-ui/icons/Timer';
import { Tooltip } from '@material-ui/core';

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
