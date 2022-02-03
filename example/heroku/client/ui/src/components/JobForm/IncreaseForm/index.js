import { useState } from "react";
import { Form, ButtonToolbar, Button } from "rsuite";
import { createJob } from "../../../api/jobs";
export default function IncreaseForm(props) {
  const [formValue, setFormValue] = useState({ queue: "increase" });

  function onChange(vals) {
    setFormValue(vals);
  }

  function onSubmit() {
    createJob(formValue);
  }

  return (
    <Form onChange={(v) => onChange({ ...formValue, ...v })}>
      <h3>Create increase job</h3>
      <Form.Group controlId="id" style={{ marginTop: 40 }}>
        <Form.ControlLabel>ID</Form.ControlLabel>
        <Form.Control name="id" autoComplete="off" />
      </Form.Group>
      <Form.Group controlId="number">
        <Form.ControlLabel>Number</Form.ControlLabel>
        <Form.Control name="number" type="number" autoComplete="off" />
      </Form.Group>
      <Form.Group>
        <ButtonToolbar>
          <Button appearance="primary" onClick={() => onSubmit()}>
            Submit
          </Button>
        </ButtonToolbar>
      </Form.Group>
    </Form>
  );
}
