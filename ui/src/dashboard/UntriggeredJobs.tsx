import * as React from 'react';
import { FC } from 'react';
import Icon from '@mui/icons-material/NewReleases';

import CardWithIcon from './CardWithIcon';

interface Props {
    value?: string;
}

const UntriggeredJobs: FC<Props> = ({ value }) => {
    return (
        <CardWithIcon
            to='/jobs?filter={"status":"untriggered"}'
            icon={Icon}
            title='Untriggered Jobs'
            subtitle={value}
        />
    );
};

export default UntriggeredJobs;
