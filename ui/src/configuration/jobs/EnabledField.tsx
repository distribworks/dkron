import * as React from "react";
import SuccessIcon from '@mui/icons-material/CheckCircle';
import FailedIcon from '@mui/icons-material/Cancel';
import { Tooltip } from '@mui/material';

const EnabledField = (props: any) => {
    if (props.record[props.source] === true) {
        return <Tooltip title="Disabled"><FailedIcon htmlColor="red" /></Tooltip>
    } else {
        return <Tooltip title="Enabled"><SuccessIcon htmlColor="green" /></Tooltip>
    } 
};

export default EnabledField;
