import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import { Title } from 'react-admin';
import { makeStyles } from '@mui/styles';
import { changeTheme } from './actions';

const useStyles = makeStyles({
    label: { width: '10em', display: 'inline-block' },
    button: { margin: '1em' },
});

const Settings = () => {
    const classes = useStyles();
    return (
        <Card>
            <Title title='Settings' />
            <CardContent>
            </CardContent>
        </Card>
    );
};

export default Settings;
