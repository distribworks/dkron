import * as React from 'react';
import { Route } from 'react-router-dom';
import Configuration from './configuration/Configuration';

const Routes = [
    <Route exact path="/configuration" render={() => <Configuration />} />,
];

export default Routes;
