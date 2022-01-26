import React, { useState, useEffect } from "react";

import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import Button from "@mui/material/Button";
import Link from "@mui/material/Link";
import DirectionsRunIcon from "@mui/icons-material/DirectionsRun";
import AlarmOffIcon from "@mui/icons-material/AlarmOff";

import { fetchJobs, retryJob, cancelJob } from "../../api/jobs/";
import ColoredStatus from "./ColoredStatus";
import { Link as RouterLink } from "react-router-dom";
import { Tooltip } from "@mui/material";

function JobsTable(props) {
  const [jobs, setJobs] = useState([]);
  const [pollCount, setPollCount] = useState(0);

  const refreshJobs = async () => {
    const resp = await fetchJobs({ limit: 100, page: 1, search: props.search });
    setJobs(resp.jobs);
  };

  useEffect(() => {
    refreshJobs();

    const interval = setInterval(() => {
      refreshJobs();
      setPollCount(pollCount + 1);
    }, 5000);

    return () => clearInterval(interval);
  }, [props.search]);

  const retryClick = async (id) => {
    await retryJob(id);
    await refreshJobs();
  };

  const cancelClick = async (id) => {
    await cancelJob(id);
    await refreshJobs();
  };

  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} aria-label="jobs table">
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell align="right">Queue</TableCell>
            <TableCell align="right">Status</TableCell>
            <TableCell align="right">Scheduled at</TableCell>
            <TableCell align="right">Last Error</TableCell>
            <TableCell align="right">Arguments</TableCell>
            <TableCell align="right">Result</TableCell>
            <TableCell align="right">Actions</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {jobs.map((row) => (
            <TableRow
              key={row.id}
              sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
            >
              <TableCell component="th" scope="row">
                <Link
                  component={RouterLink}
                  to={`/jobs/${row.id}`}
                  state={{ job: row }}
                >
                  {row.id}
                </Link>
              </TableCell>
              <TableCell align="right">{row.queue}</TableCell>
              <TableCell align="right">
                <ColoredStatus status={row.status} />
              </TableCell>
              <TableCell align="right">
                {new Date(row.scheduled_at).toLocaleString()}
              </TableCell>
              <TableCell align="right">{row.last_error}</TableCell>
              <TableCell align="right">
                {JSON.stringify(row.arguments)}
              </TableCell>
              <TableCell align="right">{JSON.stringify(row.result)}</TableCell>
              <TableCell align="right">
                {row.status !== "initialized" && (
                  <Tooltip title="Run now">
                    <Button
                      variant="contained"
                      color="primary"
                      onClick={() => retryClick(row.id)}
                    >
                      <DirectionsRunIcon />
                    </Button>
                  </Tooltip>
                )}
                {row.status === "scheduled" && (
                  <Tooltip title="Cancel">
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={() => cancelClick(row.id)}
                    >
                      <AlarmOffIcon />
                    </Button>
                  </Tooltip>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default JobsTable;
