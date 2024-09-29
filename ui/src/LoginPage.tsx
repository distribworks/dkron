import { useState } from 'react';
import { useLogin, useNotify } from 'react-admin';

const LoginPage = ({ theme }) => {
    const [token, setToken] = useState('');
    const login = useLogin();
    const notify = useNotify();

    const handleSubmit = e => {
        e.preventDefault();
        login({ token }).catch(() =>
            notify('Invalid token')
        );
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                name="token"
                type="text"
                value={token}
                onChange={e => setToken(e.target.value)}
            />
        </form>
    );
};

export default LoginPage;

