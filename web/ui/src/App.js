import './App.css';

import React, { useState } from 'react';
import Container from '@mui/material/Container';
import { Box } from '@mui/material';
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";

import SearchAppBar from './components/SearchAppBar';
import ListJobs from './routes/ListJobs';
import ViewJob from './routes/ViewJob';

function App() {

  const [search, setSearch] = useState("");


  const onSearch = (s)=>{
    setSearch(s)
  }

  return (
  <BrowserRouter>
    <SearchAppBar onSearch={ onSearch }/>

    <Box sx={{p: 2}} />

    <Container maxWidth="xl">
      <Routes>
        <Route path="/" element={ <ListJobs search={search}/>} />
        <Route path="/jobs/:id" element={ <ViewJob/>} />
      </Routes>
    </Container>
  </BrowserRouter>
  )
}

export default App;