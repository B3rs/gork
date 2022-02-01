import JobsTable from "../components/JobsTable";
import StatsWidget from "../components/StatsWidget";
import { Row, Col, Grid } from "rsuite";

function ListJobs(props) {
  return (
    <Grid>
      <Row>
        <Col md={6} sm={24}>
          <StatsWidget />
        </Col>
      </Row>
      <Row style={{ marginTop: 20 }}>
        <Col sm={24}>
          <JobsTable {...props} />
        </Col>
      </Row>
    </Grid>
  );
}

export default ListJobs;
