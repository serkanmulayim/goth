import React from 'react';

import {Route, withRouter} from "react-router-dom";

class PrivateRoute extends Route{
  constructor(props){
    super(props);
    this.state = {};
    this.isAuthenticated = false;
  }
  render() {
      console.log("from private routehistory");
      console.log(this.history);
      console.log("from props.history");
      console.log(this.props.history);
      let route;

      let ps = Object.assign({}, this.props);
      delete ps.component;
      if(localStorage.getItem("auth")) {
        let Comp = this.props.component;

        route = <Comp/>

      } else {
        //route = <Redirect to='/login' />
        this.props.history.push("/login");
      }

      return (<Route {...ps} render ={(props) => route}/>);
  }


  componentDidUpdate(){
    console.log("app component did update");
  }
}

export default PrivateRoute;
