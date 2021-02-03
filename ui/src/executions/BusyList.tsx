import * as React from "react";
import { List, Datagrid, TextField, DateField } from 'react-admin';

export const OutputPanel = ({ id, record, resource }: any) => {
    let o = record.output || "Empty output";
    return (<div className="execution-output" dangerouslySetInnerHTML={{ __html: o }} />);
};

export const BusyList = (props: any) => (
    <List {...props} pagination={ false }>
        <Datagrid rowClick="expand" isRowSelectable={ record => false } expand={<OutputPanel />}>
            <TextField source="id" sortable={false} />
            <TextField source="job_name" sortable={false} />
            <TextField source="node_name" sortable={false} />
            <DateField source="started_at" sortable={false} showTime />
        </Datagrid>
    </List>
);
export default BusyList;
