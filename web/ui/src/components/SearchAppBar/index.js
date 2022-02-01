import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { Navbar, Nav } from "rsuite";
import Input from "rsuite/Input";
import InputGroup from "rsuite/InputGroup";
import SearchIcon from "@rsuite/icons/Search";
import { Link } from "react-router-dom";

export default function SearchAppBar(props) {
  const [search, setSearch] = useState("");
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    const queryStringSearch = searchParams.get("search") || "";
    if (queryStringSearch !== "") {
      onSearch(queryStringSearch);
    }
  }, []);

  const onSearch = (search) => {
    if (search === "") {
      setSearchParams({});
    } else {
      setSearchParams({ search });
    }
    setSearch(search);
    props.onSearch(search);
  };

  return (
    <Navbar appearance="inverse">
      <Navbar.Brand href="#">
        <Link to={"/"} style={{ color: "white" }}>
          <b>HOME</b>
        </Link>
      </Navbar.Brand>
      <Nav></Nav>
      <Nav pullRight>
        <Nav.Item>
          <InputGroup>
            <Input value={search} onChange={(text) => onSearch(text)} />
            <InputGroup.Addon>
              <SearchIcon />
            </InputGroup.Addon>
          </InputGroup>
        </Nav.Item>
      </Nav>
    </Navbar>
  );
}
