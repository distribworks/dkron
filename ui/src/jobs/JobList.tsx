import {
    Datagrid,
    TextField,
    NumberField,
    DateField,
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
import EnabledField from "./EnabledField"
import { styled } from '@mui/material/styles';

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

const JobBulkActionButtons = () => (
    <Fragment>
        <BulkRunButton />
        <BulkToggleButton />
        <BulkDeleteButton />
    </Fragment>
);

const JobPagination = (props: any) => <Pagination rowsPerPageOptions={[5, 10, 25, 50, 100]} {...props} />;

const PREFIX = 'JobList';

const classes = {
    hiddenOnSmallScreens: `${PREFIX}-hiddenOnSmallScreens`,
    cell: `${PREFIX}-cell`,
};

const StyledDatagrid = styled(Datagrid)(({ theme }) => ({
    [`& .${classes.hiddenOnSmallScreens}`]: {
        display: 'table-cell',
        [theme.breakpoints.down('md')]: {
            display: 'none',
        },
    },
    [`& .${classes.cell}`]: {
        padding: "6px 8px 6px 8px",
    },
}));

const JobList = (props: any) => {
    return (
        <List {...props} filters={<JobFilter />} pagination={<JobPagination />}>
            <StyledDatagrid rowClick="show" bulkActionButtons={<JobBulkActionButtons />}>
                <TextField source="id" />
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
                <EnabledField label="Enabled" />
                <NumberField source="retries" sortable={false} />
                <StatusField />
                <DateField source="next" showTime />
                <EditButton/>
            </StyledDatagrid>
        </List>
    );
};

export default JobList;
