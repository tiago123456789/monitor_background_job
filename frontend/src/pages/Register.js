import React, { useEffect, useState } from "react"
import { Alert, Button, Container, Form, FormGroup, Input, Label } from "reactstrap"
import AuthService from "../services/AuthService"

const authService = new AuthService()

export default (props) => {
    const [credential, setCredential] = useState({ name: "", email: "", password: "" })
    const [error, setError] = useState(null);

    const handlerInputValue = (key, value) => {
        setCredential({...credential, [key]: value });
    }

    const submit = async (event) => {
        event.preventDefault()
        try {
            await authService.create(credential)
            props.history.push("/login")
        } catch(error) {
            setError(error.response.data.message)
        }
    }

    return (
        <>
            <Container>
                <h1>Register</h1>
                { error && <Alert color="danger" >{error}</Alert> }
                <Form onSubmit={submit}>
                <FormGroup>
                        <Label for="name">Name:</Label>
                        <Input
                            value={credential.name}
                            onChange={(event) => handlerInputValue("name", event.target.value)}
                            type="text" name="name"
                            id="name" />
                    </FormGroup>
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
                    <Button className="mt-1">Save</Button>
                </Form>
            </Container>
        </>
    )
}