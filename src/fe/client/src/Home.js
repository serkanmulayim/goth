import React from 'react';
import {Button} from 'react-bootstrap';
import Transport from './Transport';
import {withRouter} from "react-router-dom";

class Home extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};

  }

  handleLogout = (e) => {
    //e.preventDefault();

    Transport.DoLogoutRequest().then(resp => {
      localStorage.setItem("auth", false);
      this.props.history.push("/login");


    });
  }

  componentDidMount() {
    console.log("home did mount", this.props.history);
    Transport.DoCheckAuthRequest()
    .then(resp => {

      if(resp.success){
        this.setState({user:resp.user});
      } else {
        localStorage.removeItem("auth");
        this.props.history.push("/login")
      }
    });
  }

  render() {
    console.log(this.state);
    if(!this.state.user){
      return <div> waitiiiing</div>
    }

    return (
      <div className="Home">
        <header className="GoTH Authorization Server">
          <p>
            Welcome {this.state.user.username};
          </p>

          <Button block size="large" type="submit" onClick={this.handleLogout}>Logout</Button>
        </header>
      </div>
    );
  }
}

export default withRouter(Home);
