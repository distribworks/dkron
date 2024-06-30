import * as React from 'react';
import { FC } from 'react';
import Icon from '@mui/icons-material/ThumbUp';

import CardWithIcon from './CardWithIcon';

interface Props {
    value?: string;
}

const SuccessfulJobs: FC<Props> = ({ value }) => {
    return (
        <CardWithIcon
            to='/jobs?filter={"status":"success"}'
            icon={Icon}
            title='Successful Jobs'
            subtitle={value}
        />
    );
};

export default SuccessfulJobs;
