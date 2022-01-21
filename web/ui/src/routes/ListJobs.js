import JobsTable from '../components/JobsTable';
import { useNavigate } from "react-router-dom";

function ListJobs(props) {

  const navigate = useNavigate();

  const onJobClick = (job)=>{
    navigate(`/jobs/${job.id}`, {state: {job}});
  }
  
  return (
    <JobsTable { ...{ onJobClick, ... props }} />
  )
}

export default ListJobs;