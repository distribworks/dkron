import * as React from 'react';
import { FC } from 'react';
import Icon from '@mui/icons-material/ThumbDown';

import CardWithIcon from './CardWithIcon';

interface Props {
    value?: string;
}

const FailedJobs: FC<Props> = ({ value }) => {
    return (
        <CardWithIcon
            to='/jobs?filter={"status":"failed"}'
            icon={Icon}
            title='Failed Jobs'
            subtitle={value}
        />
    );
};

export default FailedJobs;
