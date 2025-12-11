import { HtmlHTMLAttributes, ReactNode } from 'react';
import { useState } from 'react';
import {
    Form,
    useLogin,
    useNotify,
    useSafeSetState,
} from 'react-admin';
import { styled } from '@mui/material/styles';
import { 
    Button,
    CardContent,
    CircularProgress,
    Avatar,
    Card,
    SxProps,
    TextField,
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';

const LoginPage = (props: LoginFormProps) => {
    const [token, setToken] = useState('');
    const login = useLogin();
    const notify = useNotify();
    const avatarIcon = <LockIcon />;
    const { className } = props;
    const [loading, setLoading] = useSafeSetState(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        login({ token }).catch(() => {
            setLoading(false);
            notify('Invalid token');
        });
    };

    return (
        <Root>
            <Card className={LoginClasses.card}>
                <div className={LoginClasses.avatar}>
                    <Avatar className={LoginClasses.icon}>{avatarIcon}</Avatar>
                </div>
                <form onSubmit={handleSubmit} className={className}>
                    <CardContent className={LoginFormClasses.content}>
                        <TextField
                            name="token"
                            label="Token"
                            type="text"
                            value={token}
                            onChange={e => setToken(e.target.value)}
                            fullWidth
                            variant="outlined"
                            margin="normal"
                        />

                        <Button
                            variant="contained"
                            type="submit"
                            color="primary"
                            disabled={loading}
                            fullWidth
                            className={LoginFormClasses.button}
                        >
                            {loading ? (
                                <CircularProgress
                                    className={LoginFormClasses.icon}
                                    size={19}
                                    thickness={3}
                                />
                            ) : (
                                "Sign in"
                            )}
                        </Button>
                    </CardContent>
                </form>
            </Card>
        </Root>
    );
};

export default LoginPage;

export interface LoginProps extends HtmlHTMLAttributes<HTMLDivElement> {
    avatarIcon?: ReactNode;
    backgroundImage?: string;
    children?: ReactNode;
    className?: string;
    sx?: SxProps;
}

const PREFIX = 'RaLogin';
export const LoginClasses = {
    card: `${PREFIX}-card`,
    avatar: `${PREFIX}-avatar`,
    icon: `${PREFIX}-icon`,
};

const Root = styled('div', {
    name: PREFIX,
    overridesResolver: (props, styles) => styles.root,
})(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    minHeight: '100vh',
    height: '1px',
    alignItems: 'center',
    justifyContent: 'flex-start',
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'cover',
    backgroundImage:
        'radial-gradient(circle at 50% 14em, #313264 0%, #00023b 60%, #00023b 100%)',

    [`& .${LoginClasses.card}`]: {
        minWidth: 300,
        marginTop: '6em',
    },
    [`& .${LoginClasses.avatar}`]: {
        margin: '1em',
        display: 'flex',
        justifyContent: 'center',
    },
    [`& .${LoginClasses.icon}`]: {
        backgroundColor: theme.palette.secondary.main,
    },
}));

const PREFIXF = 'RaLoginForm';

export const LoginFormClasses = {
    content: `${PREFIXF}-content`,
    button: `${PREFIXF}-button`,
    icon: `${PREFIXF}-icon`,
};

export interface LoginFormProps {
    redirectTo?: string;
    className?: string;
}
