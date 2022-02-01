import Tag from "rsuite/Tag";

const getColor = (status) => {
  switch (status) {
    case "failed":
      return "red";
    case "completed":
      return "green";
    case "scheduled":
      return "blue";
    default:
      return "default";
  }
};

export default function ColoredStatus(props) {
  return <Tag color={getColor(props.status)}>{props.status}</Tag>;
}
