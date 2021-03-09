import * as React from "react";
import SuccessIcon from '@material-ui/icons/CheckCircle';
import FailedIcon from '@material-ui/icons/Cancel';
import { Tooltip } from '@material-ui/core';

const StatusField = (props: any) => {
  return props.record === undefined ? null : (props.record[props.source] === 'success' ? <Tooltip title="Success"><SuccessIcon htmlColor="green" /></Tooltip> : <Tooltip title="Success"><FailedIcon htmlColor="red" /></Tooltip>);
};

export default StatusField;
