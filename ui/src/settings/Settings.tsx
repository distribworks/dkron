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

const Configuration = () => {
    const classes = useStyles();
    return (
        <Card>
            <Title title='Configuration' />
            <CardContent>
                <div className={classes.label}>
                    Theme
                </div>
                <Button
                    variant="contained"
                    className={classes.button}
                    color={theme === 'light' ? 'primary' : 'default'}
                    onClick={() => dispatch(changeTheme('light'))}
                >
                    Light
                </Button>
                <Button
                    variant="contained"
                    className={classes.button}
                    color={theme === 'dark' ? 'primary' : 'default'}
                    onClick={() => dispatch(changeTheme('dark'))}
                >
                    Dark
                </Button>
            </CardContent>
        </Card>
    );
};

export default Configuration;
