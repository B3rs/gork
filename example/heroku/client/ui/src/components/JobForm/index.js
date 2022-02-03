import { useState } from "react";
import IncreaseForm from "./IncreaseForm";
import LowerizeForm from "./LowerizeForm";
import { Popover, Dropdown, IconButton, Whisper } from "rsuite";
import PlusIcon from "@rsuite/icons/Plus";
import { Row } from "rsuite";

export default function JobForm(props) {
  const [selected, setSelected] = useState("");

  const renderMenu = ({ onClose, left, top, className }, ref) => {
    const handleSelect = (eventKey) => {
      onClose();
      setSelected(eventKey);
    };
    return (
      <Popover ref={ref} className={className} style={{ left, top }} full>
        <Dropdown.Menu onSelect={handleSelect}>
          <Dropdown.Item eventKey={"increase"}>Increase</Dropdown.Item>
          <Dropdown.Item eventKey={"lowerize"}>Lowerize</Dropdown.Item>
        </Dropdown.Menu>
      </Popover>
    );
  };

  return (
    <>
      <Row style={{ marginTop: 40 }}>
        <Whisper placement="bottomStart" trigger="click" speaker={renderMenu}>
          <IconButton appearance="primary" icon={<PlusIcon />} placement="left">
            New
          </IconButton>
        </Whisper>
      </Row>

      <Row style={{ marginTop: 40 }}>
        {selected === "increase" && <IncreaseForm />}
        {selected === "lowerize" && <LowerizeForm />}
      </Row>
    </>
  );
}
