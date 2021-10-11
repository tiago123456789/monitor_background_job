import React from "react"
import { BrowserRouter, Route, Switch, Redirect } from "react-router-dom"
import ListJob from "./pages/ListJobs"
import JobDetail from "./pages/JobDetail"
import Header from "./components/Header";

function App() {

  return (
    <>
      <Header/>
      <BrowserRouter>
        <Switch>
          <Route exact path="/jobs-monitored" component={ListJob} />
          <Route exact path="/jobs-monitored/:id" component={JobDetail} />
          <Redirect to="/jobs-monitored" />
        </Switch>
      </BrowserRouter>
    </>
  );
}

export default App;
