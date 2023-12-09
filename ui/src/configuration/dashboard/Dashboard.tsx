import * as React from "react";
import { Card, CardContent, CardHeader } from '@mui/material';
import { List, Datagrid, TextField } from 'react-admin';
import { TagsField } from '../TagsField'
import Leader from './Leader';
import FailedJobs from './FailedJobs';
import SuccessfulJobs from './SuccessfulJobs';
import UntriggeredJobs from './UntriggeredJobs';
import TotalJobs from './TotalJobs';

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

const styles = {
    flex: { display: 'flex' },
    flexColumn: { display: 'flex', flexDirection: 'column' },
    leftCol: { flex: 1, marginRight: '0.5em' },
    rightCol: { flex: 1, marginLeft: '0.5em' },
    singleCol: { marginTop: '1em', marginBottom: '1em' },
};

const Spacer = () => <span style={{ width: '1em' }} />;

const Dashboard = () => (
    <div>
        <Card>
            <CardHeader title="Welcome" />
            <CardContent>
                <div style={styles.flex}>
                    <div style={styles.leftCol}>
                        <div style={styles.flex}>
                            <Leader value={window.DKRON_LEADER || "devel"} />
                            <Spacer />
                            <TotalJobs value={window.DKRON_TOTAL_JOBS || "0"} />
                            <Spacer />
                            <SuccessfulJobs value={window.DKRON_SUCCESSFUL_JOBS || "0"} />
                            <Spacer />
                            <FailedJobs value={window.DKRON_FAILED_JOBS || "0"} />
                            <Spacer />
                            <UntriggeredJobs value={window.DKRON_UNTRIGGERED_JOBS || "0"} />
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
        <Card>
            <CardHeader title="Nodes" />
            <CardContent>
                <List {...fakeProps}>
                    <Datagrid isRowSelectable={ record => false }>
                        <TextField source="Name" sortable={false} />
                        <TextField source="Addr" sortable={false} />
                        <TextField source="Port" sortable={false} />
                        <TextField label="Status" source="statusText" sortable={false} />
                        <TagsField source="Tags" sortable={false} />
                    </Datagrid>
                </List>
            </CardContent>
        </Card>
    </div>
);
export default Dashboard;
