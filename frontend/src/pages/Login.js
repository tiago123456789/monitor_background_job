import React, { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { Alert, Button, Container, Form, FormGroup, Input, Label } from "reactstrap"
import AuthService from "../services/AuthService"

const authService = new AuthService()

export default (props) => {
    const [credential, setCredential] = useState({ email: "", password: "" })
    const [error, setError] = useState(null);

    const handlerInputValue = (key, value) => {
        setCredential({...credential, [key]: value });
    }

    const submit = async (event) => {
        event.preventDefault()
        try {
            const response = await authService.authenticate(credential)
            localStorage.setItem("accessToken", response.accessToken)
            setError(null)
            props.history.push("/jobs-monitored")
        } catch(error) {
            setError(error.response.data.message)
        }
    }

    return (
        <>
            <Container>
                <h1>Login</h1>
                { error && <Alert color="danger" >{error}</Alert> }
                <Form onSubmit={submit}>
                    <FormGroup>
                        <Label for="email">Email:</Label>
                        <Input
                            value={credential.email}
                            onChange={(event) => handlerInputValue("email", event.target.value)}
                            type="email" name="email"
                            id="email" />
                    </FormGroup>
                    <FormGroup>
                        <Label for="password">Password:</Label>
                        <Input
                            value={credential.password}
                            onChange={(event) => handlerInputValue("password", event.target.value)}
                            type="password" name="password"
                            id="password" />
                    </FormGroup>
                    <Button className="mt-1">Login</Button>&nbsp;
                    <Link to="/register">
                    <Button className="mt-1">Register</Button>

                    </Link>
                </Form>
            </Container>
        </>
    )
}