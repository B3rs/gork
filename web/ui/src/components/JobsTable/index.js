import React, { useState, useEffect } from "react";

import { fetchJobs, retryJob, cancelJob } from "../../api/jobs/";
import ColoredStatus from "./ColoredStatus";
import { Link } from "react-router-dom";

import Table from "rsuite/Table";

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
    <Table autoHeight height={420} data={jobs}>
      <Table.Column align="left" fixed flexGrow={1}>
        <Table.HeaderCell>ID</Table.HeaderCell>
        <Table.Cell>
          {(job) => (
            <Link to={`/jobs/${job.id}`} state={{ job: job }}>
              {job.id}
            </Link>
          )}
        </Table.Cell>
      </Table.Column>

      <Table.Column align="left" flexGrow={1}>
        <Table.HeaderCell>Queue</Table.HeaderCell>
        <Table.Cell dataKey="queue" />
      </Table.Column>

      <Table.Column align="left" flexGrow={1}>
        <Table.HeaderCell>Status</Table.HeaderCell>
        <Table.Cell>
          {(job) => <ColoredStatus status={job.status} />}
        </Table.Cell>
      </Table.Column>

      <Table.Column align="left" flexGrow={2}>
        <Table.HeaderCell>Scheduled at</Table.HeaderCell>
        <Table.Cell>
          {(job) => new Date(job.scheduled_at).toLocaleString()}
        </Table.Cell>
      </Table.Column>

      <Table.Column align="left" flexGrow={2}>
        <Table.HeaderCell>Last Error</Table.HeaderCell>
        <Table.Cell dataKey="last_error" />
      </Table.Column>

      <Table.Column align="left" flexGrow={2}>
        <Table.HeaderCell>Result</Table.HeaderCell>
        <Table.Cell>
          {(job) => <code>{JSON.stringify(job.result)}</code>}
        </Table.Cell>
      </Table.Column>

      <Table.Column align="right" flexGrow={1} fixed>
        <Table.HeaderCell>Actions</Table.HeaderCell>
        <Table.Cell>
          {(job) => {
            switch (job.status) {
              case "scheduled":
                return (
                  <>
                    <a href="#" onClick={() => retryClick(job.id)}>
                      Run
                    </a>{" "}
                    |
                    <a href="#" onClick={() => cancelClick(job.id)}>
                      Cancel
                    </a>
                  </>
                );
              case "initialized":
                return <></>;
              case "failed":
                return (
                  <>
                    <a href="#" onClick={() => retryClick(job.id)}>
                      Retry
                    </a>
                  </>
                );
              case "canceled":
                return (
                  <>
                    <a href="#" onClick={() => retryClick(job.id)}>
                      Run
                    </a>
                  </>
                );
              case "completed":
                return (
                  <>
                    <a href="#" onClick={() => retryClick(job.id)}>
                      Run again
                    </a>
                  </>
                );
              default:
            }
          }}
        </Table.Cell>
      </Table.Column>
    </Table>
  );
}

export default JobsTable;
