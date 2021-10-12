import React from "react";
import { Redirect, Route } from "react-router-dom";
import AuthService from "../services/AuthService";

const authService = new AuthService

export default ({ component: Component, ...rest }) => {
    return (
        <Route {...rest} render={(props) => (       
            authService.isAuthenticated()
            ? <Component {...props} />
            : <Redirect to='/login' />
         )} />
    )
} 