import { HtmlHTMLAttributes, ReactNode } from 'react';
import { useState } from 'react';
import {
    Form,
    useLogin,
    useNotify,
    TextInput,
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
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';

const LoginPage = (props: LoginFormProps) => {
    const [token, setToken] = useState('');
    const login = useLogin();
    const notify = useNotify();
    const avatarIcon = <LockIcon />;
    const { className } = props;
    const [loading, setLoading] = useSafeSetState(false);

    const handleSubmit = e => {
        e.preventDefault();
        setLoading(true);
        login({ token }).catch(() => {
            setLoading(false);
            notify('Invalid token');
        });
    };

    return (
        <Card className={LoginClasses.card}>
                <div className={LoginClasses.avatar}>
                    <Avatar className={LoginClasses.icon}>{avatarIcon}</Avatar>
                </div>
                <StyledForm
                    onSubmit={handleSubmit}
                    mode="onChange"
                    noValidate
                    className={className}
                >
                    <CardContent className={LoginFormClasses.content}>
                        <TextInput
                            name="token"
                            type="text"
                            value={token}
                            onChange={e => setToken(e.target.value)}
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
                </StyledForm>
        </Card>
        
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
        backgroundColor: theme.palette.secondary[500],
    },
}));

const PREFIXF = 'RaLoginForm';

export const LoginFormClasses = {
    content: `${PREFIXF}-content`,
    button: `${PREFIXF}-button`,
    icon: `${PREFIXF
    }-icon`,
};

const StyledForm = styled(Form, {
    name: PREFIXF,
    overridesResolver: (props, styles) => styles.root,
})(({ theme }) => ({
    [`& .${LoginFormClasses.content}`]: {
        width: 300,
    },
    [`& .${LoginFormClasses.button}`]: {
        marginTop: theme.spacing(2),
    },
    [`& .${LoginFormClasses.icon}`]: {
        margin: theme.spacing(0.3),
    },
}));

export interface LoginFormProps {
    redirectTo?: string;
    className?: string;
}
