import React from "react"
import { BrowserRouter, Route, Switch, Redirect } from "react-router-dom"
import ListJob from "./pages/ListJobs"
import JobDetail from "./pages/JobDetail"
import Header from "./components/Header";
import Login from "./pages/Login";
import PrivateRoute from "./components/PrivateRoute"
import Register from "./pages/Register";

function App() {

  return (
    <>
      <Header/>
      <BrowserRouter>
        <Switch>
          <Route exact path="/login" component={Login} />
          <Route exact path="/register" component={Register} />
          <PrivateRoute exact path="/jobs-monitored" component={ListJob} />
          <PrivateRoute exact path="/jobs-monitored/:id" component={JobDetail} />
          <Redirect to="/login" />
        </Switch>
      </BrowserRouter>
    </>
  );
}

export default App;
