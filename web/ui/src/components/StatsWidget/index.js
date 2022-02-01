import { useState, useEffect } from "react";

import { Col, Row, Panel, Grid } from "rsuite";
import { fetchStats } from "../../api/jobs";

function Card(props) {
  const { queue } = props;

  return (
    <Panel {...props} bordered header={"Queue: " + queue.name}>
      <p>
        Scheduled: {queue.scheduled} <br />
        Initialized: {queue.initialized} <br />
        Completed: {queue.completed} <br />
        Failed: {queue.failed} <br />
      </p>
    </Panel>
  );
}

function StatsWidget(props) {
  const [stats, setStats] = useState({});

  const getStats = async () => {
    const stats = await fetchStats();
    setStats(stats);
  };

  useEffect(() => {
    getStats();
  }, []);

  const queues = stats.queues || [];

  return (
    <Grid>
      <Row>
        {queues.map((q) => (
          <Col md={6} sm={24} key={"card" + q.name}>
            <Card key={q.name} queue={q} />
          </Col>
        ))}
      </Row>
    </Grid>
  );
}

export default StatsWidget;
