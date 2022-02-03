import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import { fetchJob } from "../api/jobs";
import { Grid, Col, Row, Table, Divider } from "rsuite";

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
    { label: "Arguments", value: JSON.stringify(job.arguments) },
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
    <Grid>
      <Row>
        <Col sm={24}>
          <h3>Job: {job.id}</h3>
        </Col>
      </Row>
      <Row>
        <Col sm={24}>
          <Divider />
          <Table data={fields} showHeader={false} autoHeight={true}>
            <Table.Column align="right">
              <Table.HeaderCell>Pippo</Table.HeaderCell>
              <Table.Cell dataKey="label" />
            </Table.Column>
            <Table.Column align="left" flexGrow={4}>
              <Table.HeaderCell>Pippo</Table.HeaderCell>
              <Table.Cell dataKey="value" />
            </Table.Column>
          </Table>
        </Col>
      </Row>
    </Grid>
  );
}

export default ViewJob;
