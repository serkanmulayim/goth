import React from 'react';
import Login from './Login';
import Home from './Home';
import { HashRouter as Router, Route, Redirect} from "react-router-dom";

const PrivateRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={(props) => (
    localStorage.getItem("auth")
      ? <Component {...props} />
      : <Redirect to='/login' />
  )} />
)

class App extends React.Component {
  render() {
      return (
      <div className="App">

          <Router>
          <div>
            <PrivateRoute exact path="/" component={Home} />
            <Route exact path="/login" component={Login} />
          </div>

          </Router>
      </div>
    );
  }
}

export default App;
