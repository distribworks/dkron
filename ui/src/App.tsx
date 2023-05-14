import * as React from "react";
import Dashboard from './dashboard';
import { Admin, Resource } from 'react-admin';
import dataProvider from './dataProvider';
import jobs from './jobs';
import { BusyList } from './executions/BusyList';
import PlayCircleOutlineIcon from '@material-ui/icons/PlayCircleOutline';
import { Layout } from './layout';
import customRoutes from './routes';
import themeReducer from './themeReducer';
import createAdminStore from './createAdminStore'
import { Provider } from "react-redux";
import { createHashHistory } from "history";
import { persistStore } from 'redux-persist';
import { PersistGate } from "redux-persist/integration/react";

declare global {
    interface Window {
        DKRON_API_URL: string;
        DKRON_LEADER: string;
        DKRON_UNTRIGGERED_JOBS: string;
        DKRON_FAILED_JOBS: string;
        DKRON_SUCCESSFUL_JOBS: string;
        DKRON_TOTAL_JOBS: string;
    }
}

const authProvider = () => Promise.resolve();
const history = createHashHistory();
const store = createAdminStore({authProvider, dataProvider, history});
const persistor = persistStore(store);

const App = () => (
    <Provider store={store}>
        <PersistGate loading={null} persistor={persistor}>
            <Admin
                dashboard={Dashboard} 
                authProvider={authProvider}
                dataProvider={dataProvider}
                history={history}
                layout={Layout}
                customRoutes={customRoutes}
                customReducers={{ theme: themeReducer }}
            >
                <Resource name="jobs" {...jobs} />
                <Resource name="busy" options={{ label: 'Busy' }} list={BusyList} icon={PlayCircleOutlineIcon} />
                <Resource name="executions" />
                <Resource name="members" />
            </Admin>
        </PersistGate>
    </Provider>
);

export default App;
