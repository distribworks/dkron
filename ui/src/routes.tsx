import { Route } from 'react-router-dom';
import Configuration from './settings/Settings';

const Routes = [
    <Route exact path="/configuration" render={() => <Configuration />} />,
];

export default Routes;
