import React from 'react';
import PrivateRoute from './PrivateRoute';
import Login from './Login';
import Home from './Home';
import { HashRouter as Router, Route} from "react-router-dom";


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
