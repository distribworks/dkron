import JobList from './JobList';
import { JobEdit, JobCreate } from './JobEdit';
import JobShow from './JobShow';
import JobIcon from '@mui/icons-material/Update';

const jobs = {
    list: JobList,
    edit: JobEdit,
    create: JobCreate,
    show: JobShow,
    icon: JobIcon
};
export default jobs;
