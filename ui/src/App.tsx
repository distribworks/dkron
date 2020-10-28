import * as React from "react";
import Dashboard from './dashboard';
import { Admin, Resource } from 'react-admin';
import myDataProvider from './dataProvider';
import jobs from './jobs';
import { BusyList } from './executions/BusyList';
import PlayCircleOutlineIcon from '@material-ui/icons/PlayCircleOutline';
import { Layout } from './layout';
import customRoutes from './routes';
import themeReducer from './themeReducer';

const App = () => (
    <Admin
        dashboard={Dashboard} 
        dataProvider={myDataProvider} 
        layout={Layout}
        customRoutes={customRoutes}
        customReducers={{ theme: themeReducer }}
    >
        <Resource name="jobs" {...jobs} />
        <Resource name="busy" options={{ label: 'Busy' }} list={BusyList} icon={PlayCircleOutlineIcon} />
        <Resource name="executions" />
        <Resource name="members" />
    </Admin>
);

export default App;
