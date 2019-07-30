import React from 'react';
import {Navbar, Nav, NavDropdown} from 'react-bootstrap';
import AdminApi from './api/AdminApi';
import SigningKeys from './components/SigningKeys';
import Admins from './components/Admins';
import OAuthApps from './components/OAuthApps';
import InternalUsers from './components/InternalUsers';
import Databases from './components/Databases';
import {withRouter} from "react-router-dom";

class Home extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      activeComponent: Admins
    };
  }

  handleLogout = (e) => {
    e.preventDefault();

    AdminApi.DoLogoutRequest().then(resp => {
      localStorage.removeItem("auth");
      this.props.history.push("/login");
    });
  }

  componentDidMount() {


    AdminApi.DoCheckAuthRequest()
      .then(resp => {


        if (resp.success) {
          let newState = Object.assign({}, this.state);
          newState.admin = resp.admin;
          this.setState(newState);
        } else {
          localStorage.removeItem("auth");
          this.props.history.push("/login");
        }
      });
  }

  handleLogout = () => {


    AdminApi.DoLogoutRequest().then(resp => {
      localStorage.removeItem("auth");
      this.props.history.push("/login");
    });
  }

  handleSelect = (eventKey) => {
    let Comp;
    switch(eventKey) {
      case "oauth":
        Comp = OAuthApps;
        break;
      case "signingkeys":
        Comp = SigningKeys;
        break;
      case "databases":
        Comp = Databases;
        break;
      case "internalusers":
        Comp = InternalUsers;
        break;
      case "admins":
        Comp = Admins;
        break;
      default:
        Comp = OAuthApps;
    }


    let newState = Object.assign({}, this.state);
    newState.activeComponent = Comp;
    this.setState(newState);
  }

  getUsernameString() {
    if(!this.state.admin) {
      return "Admin";
    }
    return this.state.admin.firstname + " " + this.state.admin.lastname;
  }


render() {

  return (

    <div className = "Home" >
        <header>
            <Navbar bg = "dark" variant = "dark" >
                <Navbar.Brand href = "#home" > GoTH </Navbar.Brand>
                <Nav className = "mr-auto" defaultActiveKey={"admins"} onSelect={this.handleSelect}>
                    <Nav.Item >
                        {/* <Nav.Link href = "#applications">Applications</Nav.Link> */}
                        <Nav.Link eventKey="oauth"> OAuth </Nav.Link>
                    </Nav.Item>
                    <Nav.Item >
                        <Nav.Link eventKey="signingkeys">Signature Keys</Nav.Link>
                    </Nav.Item>

                    <Nav.Item >

                        <Nav.Link eventKey = "databases">Databases</Nav.Link>
                    </Nav.Item>
                    <Nav.Item >
                        <Nav.Link eventKey = "internalusers">Internal Users</Nav.Link>

                    </Nav.Item>
                    <Nav.Item >
                        <Nav.Link eventKey = "admins">Admins</Nav.Link>
                    </Nav.Item>
                </Nav>
                {/* <Navbar.Collapse>{this.state.user}</Navbar.Collapse> */}
                <Nav>
                    <NavDropdown title={this.getUsernameString()} id="usermenu" >
                        <NavDropdown.Item onSelect={this.handleLogout}>Sign Out</NavDropdown.Item>
                        <NavDropdown.Item onSelect={this.handleLogout}>Change Password</NavDropdown.Item>
                    </NavDropdown>
                    {/* <Button variant = "outline-info" size = "sm" onClick = {this.handleLogout}> Sign out </Button> */}
                </Nav>
            </Navbar>

        </header>

        <div>
            {React.createElement(this.state.activeComponent,{user:this.state.user, handleLogout:this.handleLogout})}
        </div>
    </div>

  )};


}


export default withRouter(Home);
