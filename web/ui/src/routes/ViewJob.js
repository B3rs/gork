import { useEffect, useState } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import { fetchJob } from '../api/jobs';
import Typography from '@mui/material/Typography';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Divider from '@mui/material/Divider';
import Box from '@mui/material/Divider';


function ViewJob(props) {
  const location = useLocation();
  const { id } = useParams();
  const [ job, setJob ] = useState({});

  const loadJob = async (id) => {
    setJob(await fetchJob(id))
  }

  useEffect(() => {
    if (location.state.job !== undefined) {
      setJob(location.state.job);
      return
    }

    loadJob(id);
  },[])


  return (
    <>
      <Typography variant="h3" component="h1">
        Job: {job.id}
      </Typography>
      
      <Box sx={{p: 2}} />
      <Divider />
      <Box sx={{p: 2}} />

      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
          <TableBody>
              <TableRow key="id" sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                <TableCell>
                  ID
                </TableCell>
                <TableCell align="left">
                  {job.id}
                </TableCell>
              </TableRow>
              <TableRow key="id" sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                <TableCell>
                  Queue
                </TableCell>
                <TableCell align="left">
                  {job.queue}
                </TableCell>
              </TableRow>
          </TableBody>
        </Table>
      </TableContainer>
 
    </>
  )
}

export default ViewJob;