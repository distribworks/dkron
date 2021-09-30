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

const initialState = () => ({
    theme: localStorage.getItem("dkron-ui-theme"),
});

const App = () => (
    <Admin
        dashboard={Dashboard} 
        dataProvider={dataProvider}
        layout={Layout}
        customRoutes={customRoutes}
        initialState={initialState}
        customReducers={{ theme: themeReducer }}
    >
        <Resource name="jobs" {...jobs} />
        <Resource name="busy" options={{ label: 'Busy' }} list={BusyList} icon={PlayCircleOutlineIcon} />
        <Resource name="executions" />
        <Resource name="members" />
    </Admin>
);

export default App;
