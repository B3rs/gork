import "./App.css";

import React, { useState } from "react";
import "rsuite/dist/rsuite.min.css";
import { HashRouter, Routes, Route } from "react-router-dom";
import SearchAppBar from "./components/SearchAppBar";
import ListJobs from "./routes/ListJobs";
import ViewJob from "./routes/ViewJob";
import { Container, Content, Footer, Header } from "rsuite";

function App() {
  const [search, setSearch] = useState("");

  const onSearch = (s) => {
    setSearch(s);
  };

  return (
    <HashRouter>
      <Container>
        <Header>
          <SearchAppBar onSearch={onSearch} />
        </Header>
        <Content style={{ marginTop: 20 }}>
          <Routes>
            <Route path="/" element={<ListJobs search={search} />} />
            <Route path="/jobs/:id" element={<ViewJob />} />
          </Routes>
        </Content>
        <Footer> </Footer>
      </Container>
    </HashRouter>
  );
}

export default App;
