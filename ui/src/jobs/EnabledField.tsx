import * as React from "react";
import SuccessIcon from '@material-ui/icons/CheckCircle';
import FailedIcon from '@material-ui/icons/Cancel';
import { Tooltip } from '@material-ui/core';

const EnabledField = (props: any) => {
    if (props.record[props.source] === true) {
        return <Tooltip title="Disabled"><FailedIcon htmlColor="red" /></Tooltip>
    } else {
        return <Tooltip title="Enabled"><SuccessIcon htmlColor="green" /></Tooltip>
    } 
};

export default EnabledField;
