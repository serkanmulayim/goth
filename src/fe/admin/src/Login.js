import React from 'react';
import logo from './logo.svg';
import './App.css'
import {Form, Button} from "react-bootstrap";
import AdminApi from './api/AdminApi';
import {withRouter} from "react-router-dom";

class Login extends React.Component {
  constructor(props) {
    super(props);
    //states username and message
    this.state={message:""};
  }

  submitHandler = (event) => {
    event.preventDefault();

    if(!event.target.username.value || !event.target.password.value){
      this.setState({message:"Please enter username and password"});
    } else {
      AdminApi.DoLoginRequest(event.target.username.value,
            event.target.password.value)
      .then(resp => {
        if(!resp.success) {
          localStorage.setItem("auth", false);

          let msg;
          if(resp.status === 401){
            msg = "Wrong username or password";
          } else {
            msg = "An error occured, please try again later";
          }
          this.setState({message:msg});
        } else { //success
          localStorage.setItem("auth", true);
          this.props.history.push("/");
        }

    });
    }
  }


  componentDidMount() {
    AdminApi.DoCheckAuthRequest()
    .then(resp => {
      if(resp.success){
        localStorage.setItem("auth", true);
        this.props.history.push("/");

      } else {
        localStorage.removeItem("auth");
      }
    });
  }


  render() {

    return (

      <div className="Login" style={{position:'relative',left:'30%'}}>
          <header className="GoTH Authorization Server"/>


          <div style={{width:'300px', height:'380px', textAlign:'center', border:'2px solid #000099', padding:'20px', position:'center', margin:'100px'}}>

              <Form onSubmit={this.submitHandler}>
                  <img src={logo} className="App-logo" alt="logo"/>
                  <br/><br/>Admin Login<br/>

                  <Form.Group controlId="username" size="large">
                      {/* <Form.Label>Username</Form.Label> */}
                      <Form.Control autoFocus type="text" placeholder="Username" />
                  </Form.Group>
                  <Form.Group controlId="password" size="large">
                      {/* <Form.Label>Password</Form.Label> */}
                      <Form.Control type="password" placeholder="Password"/>
                  </Form.Group>
                  <Button block size="large" type="submit">Login</Button>
                  <Form.Label style={{color:'red'}}>{this.state.message}</Form.Label>
              </Form>
          </div>

      </div>

    );
  }
}

export default withRouter(Login);
