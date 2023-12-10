import Dashboard from './dashboard';
import { Admin, Resource, CustomRoutes } from 'react-admin';
import { Route } from "react-router-dom";
import dataProvider from './dataProvider';
import jobs from './jobs';
import { BusyList } from './executions/BusyList';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import { Layout } from './layout';
import { createHashHistory } from "history";

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

export const App = () => <Admin
    dashboard={Dashboard}
    authProvider={authProvider}
    dataProvider={dataProvider}
    layout={Layout}
>
    <Resource name="jobs" {...jobs} />
    <Resource name="busy" options={{ label: 'Busy' }} list={BusyList} icon={PlayCircleOutlineIcon} />
    <Resource name="executions" />
    <Resource name="members" />
    <CustomRoutes>
        <Route path="/settings" element={<Settings />} />
    </CustomRoutes>
</Admin>;
