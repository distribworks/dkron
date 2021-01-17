import * as React from "react";
import { makeStyles } from '@material-ui/core/styles';
import SuccessIcon from '@material-ui/icons/CheckCircle';
import FailedIcon from '@material-ui/icons/Cancel';

const StatusField = (props: any) => {
  return (props.record[props.source] === 'success' ? <SuccessIcon htmlColor="green" /> : <FailedIcon htmlColor="red" />);
};

export default StatusField;
