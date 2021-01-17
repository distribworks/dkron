import * as React from "react";
import { DateField, DateFieldProps } from 'react-admin';

export const ZeroDateField: React.FC<DateFieldProps> = (props) => {
    if (props.record !== undefined && props.source !== undefined) {
        if (props.record[props.source] === "0001-01-01T00:00:00Z") {
            props.record[props.source] = null;
        }
    }
    return (<DateField {...props}/>);
}

export default ZeroDateField;
