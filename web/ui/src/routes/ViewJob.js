import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import { fetchJob } from "../api/jobs";
import Typography from "@mui/material/Typography";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell, { tableCellClasses } from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import Divider from "@mui/material/Divider";
import Box from "@mui/material/Divider";
import { styled } from "@mui/material/styles";

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: theme.palette.common.black,
    color: theme.palette.common.white,
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
  },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  "&:nth-of-type(odd)": {
    backgroundColor: theme.palette.action.hover,
  },
  // hide last border
  "&:last-child td, &:last-child th": {
    border: 0,
  },
}));

function ViewJob(props) {
  const location = useLocation();
  const { id } = useParams();
  const [job, setJob] = useState({});

  const loadJob = async (id) => {
    setJob(await fetchJob(id));
  };

  useEffect(() => {
    if (location.state.job !== undefined) {
      setJob(location.state.job);
      return;
    }

    loadJob(id);
  }, []);

  const fields = [
    { label: "ID", value: job.id },
    { label: "Queue", value: job.queue },
    { label: "Status", value: JSON.stringify(job.status) },
    { label: "Result", value: JSON.stringify(job.result) },
    { label: "Last Error", value: job.last_error },
    { label: "Retry Count", value: job.retry_count },
    { label: "Options", value: JSON.stringify(job.options) },
    {
      label: "Scheduled at",
      value: new Date(job.scheduled_at).toLocaleString(),
    },
    { label: "Started at", value: new Date(job.started_at).toLocaleString() },
    { label: "Created at", value: new Date(job.created_at).toLocaleString() },
    { label: "Updated at", value: new Date(job.updated_at).toLocaleString() },
  ];

  return (
    <>
      <Typography variant="h3" component="h1">
        Job: {job.id}
      </Typography>

      <Box sx={{ p: 2 }} />

      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
          <TableBody>
            {fields.map((field) => (
              <StyledTableRow key={field.label}>
                <StyledTableCell>{field.label}</StyledTableCell>
                <StyledTableCell>{field.value}</StyledTableCell>
              </StyledTableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
}

export default ViewJob;
