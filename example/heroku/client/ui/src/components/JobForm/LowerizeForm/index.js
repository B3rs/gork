import { useState } from "react";
import { Form, ButtonToolbar, Button, DatePicker } from "rsuite";
import { createJob } from "../../../api/jobs";

export default function LowerizeForm(props) {
  const [formValue, setFormValue] = useState({ queue: "lowerize" });
  const [response, setResponse] = useState({});

  function onChange(vals) {
    setFormValue(vals);
  }

  function onSelectedDate(date) {
    setFormValue({ ...formValue, scheduled_at: date });
  }

  async function onSubmit() {
    const job = await createJob(formValue);
    setResponse(job);
  }

  return (
    <>
      <h3>Create lowerize job</h3>
      <Form onChange={(v) => onChange({ ...formValue, ...v })}>
        <Form.Group controlId="id" style={{ marginTop: 40 }}>
          <Form.ControlLabel>ID</Form.ControlLabel>
          <Form.Control name="id" autoComplete="off" />
        </Form.Group>
        <Form.Group controlId="string">
          <Form.ControlLabel>String</Form.ControlLabel>
          <Form.Control name="string" autoComplete="off" />
        </Form.Group>
        <Form.Group>
          <DatePicker
            format="yyyy-MM-dd HH:mm:ss"
            ranges={[
              {
                label: "Now",
                value: new Date(),
              },
            ]}
            onSelect={(date) => onSelectedDate(date)}
            style={{ width: 260 }}
          />
        </Form.Group>
        <Form.Group>
          <ButtonToolbar>
            <Button appearance="primary" onClick={() => onSubmit()}>
              Submit
            </Button>
          </ButtonToolbar>
        </Form.Group>
      </Form>
      <p style={{ marginTop: 40 }}>
        <code>{JSON.stringify(response)}</code>
      </p>
    </>
  );
}
