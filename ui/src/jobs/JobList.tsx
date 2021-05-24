import * as React from "react";
import {
    Datagrid,
    TextField,
    NumberField,
    DateField,
    BooleanField,
    EditButton,
    Filter,
    TextInput,
    List,
    SelectInput,
    BulkDeleteButton,
    BooleanInput,
    Pagination
} from 'react-admin';
import { Fragment } from 'react';
import BulkRunButton from "./BulkRunButton"
import BulkToggleButton from "./BulkToggleButton"
import StatusField from "./StatusField"
import { makeStyles } from '@material-ui/core/styles';

const JobFilter = (props: any) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <SelectInput source="status" choices={[
            { id: 'success', name: 'Success' },
            { id: 'failed', name: 'Failed' },
            { id: 'untriggered', name: 'Waiting to Run' },
        ]} />
        <BooleanInput source="disabled"/>
    </Filter>
);

const JobBulkActionButtons = (props: any) => (
    <Fragment>
        <BulkRunButton {...props} />
        <BulkToggleButton {...props} />
        <BulkDeleteButton {...props} />
    </Fragment>
);

const JobPagination = (props: any) => <Pagination rowsPerPageOptions={[5, 10, 25, 50, 100]} {...props} />;

const useStyles = makeStyles(theme => ({
    hiddenOnSmallScreens: {
        display: 'table-cell',
        [theme.breakpoints.down('md')]: {
            display: 'none',
        },
    },
    cell: {
        padding: "6px 8px 6px 8px",
    },
}));

const JobList = (props: any) => {
    const classes = useStyles();
    return (
        <List {...props} filters={<JobFilter />} bulkActionButtons={<JobBulkActionButtons />} pagination={<JobPagination />}>
            <Datagrid rowClick="show" classes={{rowCell: classes.cell}}>
                <TextField source="id"
                    cellClassName={classes.hiddenOnSmallScreens}
                    headerClassName={classes.hiddenOnSmallScreens} />
                <TextField source="displayname" label="Display name" />
                <TextField source="timezone" sortable={false}
                    cellClassName={classes.hiddenOnSmallScreens}
                    headerClassName={classes.hiddenOnSmallScreens} />
                <TextField source="schedule" />
                <NumberField source="success_count" 
                    cellClassName={classes.hiddenOnSmallScreens}
                    headerClassName={classes.hiddenOnSmallScreens} />
                <NumberField source="error_count" 
                    cellClassName={classes.hiddenOnSmallScreens}
                    headerClassName={classes.hiddenOnSmallScreens} />
                <DateField source="last_success" showTime />
                <DateField source="last_error" showTime />
                <BooleanField source="disabled" sortable={false} />
                <NumberField source="retries" sortable={false} />
                <StatusField source="status" sortable={false} />
                <DateField source="next" showTime />
                <EditButton/>
            </Datagrid>
        </List>
    );
};

export default JobList;
