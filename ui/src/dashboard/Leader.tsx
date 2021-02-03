import * as React from 'react';
import { FC } from 'react';
import Icon from '@material-ui/icons/DeviceHub';

import CardWithIcon from './CardWithIcon';

interface Props {
    value?: string;
}

const Leader: FC<Props> = ({ value }) => {
    return (
        <CardWithIcon
            to="/jobs"
            icon={Icon}
            title='Leader'
            subtitle={value}
        />
    );
};

export default Leader;
