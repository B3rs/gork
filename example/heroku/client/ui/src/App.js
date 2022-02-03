import "./App.css";

import "rsuite/dist/rsuite.min.css";
import JobForm from "./components/JobForm";
import { Grid } from "rsuite";

function App() {
  return (
    <Grid>
      <JobForm />
    </Grid>
  );
}

export default App;
