import { AuthProvider, UserIdentity } from 'react-admin';

const authProvider: AuthProvider = {
    login: ({ token }) => {
        localStorage.setItem('token', token);
        return Promise.resolve();
    },
    logout: () => {
        localStorage.removeItem('token');
        return Promise.resolve();
    },
    checkAuth: () =>
        localStorage.getItem('token') ? Promise.resolve() : Promise.reject(),
    checkError: (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            localStorage.removeItem('token');
            return Promise.reject();
        }
        // other error code (404, 500, etc): no need to log out
        return Promise.resolve();
    },
    getIdentity: () => Promise.resolve({ id: 'user', fullName: 'User' } as UserIdentity),
    getPermissions: () => Promise.resolve(),
};

export default authProvider;
