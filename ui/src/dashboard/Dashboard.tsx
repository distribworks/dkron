import * as React from "react";
import { Card, CardContent, CardHeader } from '@material-ui/core';
import { List, Datagrid, TextField } from 'react-admin';
import { TagsField } from '../TagsField'

let fakeProps = {
    basePath: "/members",
    count: 10,
    hasCreate: false,
    hasEdit: false,
    hasList: true,
    hasShow: false,
    location: { pathname: "/", search: "", hash: "", state: undefined },
    match: { path: "/", url: "/", isExact: true, params: {} },
    options: {},
    permissions: null,
    resource: "members"
}

const Dashboard = () => (
    <Card>
        <CardHeader title="Nodes" />
        <CardContent>
            <List {...fakeProps}>
                <Datagrid isRowSelectable={ record => false }>
                    <TextField source="Name" sortable={false} />
                    <TextField source="Addr" sortable={false} />
                    <TextField source="Port" sortable={false} />
                    <TextField source="Status" sortable={false} />
                    <TagsField source="Tags" sortable={false} />
                </Datagrid>
            </List>
        </CardContent>
    </Card>
);
export default Dashboard;
