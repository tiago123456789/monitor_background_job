import React, { useEffect, useState } from "react"
import { Button, Container, Form, FormGroup, Input, Label, Table } from "reactstrap"
import JobService from "../services/JobService"

const jobService = new JobService()

export default (props) => {
    const [alert, setAlert] = useState({ timeInMinutes: 5 })
    const [alerts, setAlerts] = useState([])

    const handlerInputValue = (key, value) => {
        setAlert({...alert, [key]: value })
    }

    const submit = async (event) => {
        event.preventDefault();
        alert.jobId = props.match.params.id
        await jobService.createAlert(alert)
        setAlert({ 
            name: "", timeInMinutes: 5, url: "",
            payload: ""
        })
    }

    const getAlerts = async () => {
        const registers = await jobService.getAllAlerts(props.match.params.id)
        setAlerts(registers || [])
    }


    useEffect(() => {
        (async () => {
            await getAlerts()
        })()


    }, [])


    return (
        <>
            <Container>
                <h1>Alerts</h1>
                <div>
                    <Form>
                        <FormGroup>
                            <Label for="name">Name:</Label>
                            <Input
                                value={alert.name}
                                onChange={(event) => handlerInputValue("name", event.target.value)}
                                type="text" name="name"
                                id="name" />
                        </FormGroup>
                        <FormGroup>
                            <Label for="timeInMinutes">Time in minutes:</Label>
                            <Input
                                value={alert.timeInMinutes}
                                onChange={(event) => handlerInputValue("timeInMinutes", event.target.value)}
                                type="number" name="timeInMinutes"
                                id="timeInMinutes" min="5" />
                        </FormGroup>
                        <FormGroup>
                            <Label for="job">Url:</Label>
                            <Input
                                value={alert.url}
                                onChange={(event) => handlerInputValue("url", event.target.value)}
                                type="url" name="url"
                                id="url" />
                        </FormGroup>
                        <FormGroup>
                            <Label for="job">Payload:</Label>
                            <Input type="textarea"
                                value={alert.payload}
                                onChange={(event) => handlerInputValue("payload", event.target.value)}
                                name="job"
                                id="job" placeholder="Type name the job" />
                        </FormGroup>
                        <Button type="submit" onClick={submit} className="mt-1">Novo job</Button>
                    </Form>
                </div>
                <br/>
                <h3>Alerts</h3>
                <Table>
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Name</th>
                            <th>Time in minutes</th>
                            <th>Url</th>
                            <th>Payload</th>
                        </tr>
                    </thead>
                    <tbody>
                        { alerts.length > 0 &&
                            (alerts).map(item => {
                                return (
                                    <tr key={item.id}>
                                        <th scope="row">{item.id}</th>
                                        <td>{item.name || "" }</td>
                                        <td>{item.timeInMinutes}</td>
                                        <td>{item.url}</td>
                                        <td>{item.payload}</td>

                                    </tr>
                                )
                            })
                        }

                    </tbody>
                </Table>
            </Container>
        </>
    )
}