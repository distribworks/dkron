import JobList from './JobList';
import { JobEdit, JobCreate } from './JobEdit';
import JobShow from './JobShow';
import JobIcon from '@material-ui/icons/Update';

const jobs = {
    list: JobList,
    edit: JobEdit,
    create: JobCreate,
    show: JobShow,
    icon: JobIcon
};
export default jobs;
