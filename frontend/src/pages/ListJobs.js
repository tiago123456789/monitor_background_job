import React, { useState, useEffect } from "react"
import { Button, Container, Form, FormGroup, Input, Label, Table } from "reactstrap"
import moment from "moment"
import { Link } from "react-router-dom"
import {CopyToClipboard} from 'react-copy-to-clipboard';
import JobService from "../services/JobService"

const jobService = new JobService()

export default (props) => {
    const [newJob, setNewJob] = useState({})
    const [jobs, setJobs] = useState([])
    const [eventNotifications, setEventNotifications] = useState([])

    const handlerInputValue = (key, value) => {
        setNewJob({ [key]: value })
    }

    const getJobs = async () => {
        const jobs = await jobService.getAll()
        setJobs(jobs)
    }

    const getEventNotifications = async () => {
        const registers = await jobService.getEventNotifications()
        setEventNotifications(registers)
    }   

    const generateLinkJob = (job) => {
        return jobService.generateLinkJob(job)
    }

    const submit = async (event) => {
        event.preventDefault();
        await jobService.create({ name: newJob.name })
        setNewJob({ name: "" })
        await getJobs();
    }

    const getLastNotication = (job) => {
        if (eventNotifications[job.id] != null) {
            return moment(eventNotifications[job.id]["createdAt"])
                .format("DD/MM/YY HH:mm:ss")
        }
        return "Don't have recently. Check your job"
    }

    useEffect(() => {
        (async () => {
            await getJobs()
            await getEventNotifications()
        })()

        const intervalLastNotification = setInterval(async () => {
            await getEventNotifications()
        }, 30000)

        return () => clearInterval(intervalLastNotification);
    }, [])

    return (
        <>
            <Container>
                <h1 className="text-center">Jobs</h1>

                <div>
                    <Form>
                        <FormGroup>
                            <Label for="job">Name:</Label>
                            <Input 
                            value={newJob.name}
                            onChange={(event) => handlerInputValue("name", event.target.value)}
                            type="text" name="job" 
                            id="job" placeholder="Type name the job" />
                        </FormGroup>
                        <Button type="submit" onClick={submit} className="mt-1">Novo job</Button>
                    </Form>
                </div>
                <br/>
                <Table>
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Job</th>
                            <th>Most recently notification</th>
                            <th>Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            jobs.map(job => {
                                return (
                                    <tr key={job.id}>
                                        <th scope="row">{job.id}</th>
                                        <td>{job.Name}</td>
                                        <td>{getLastNotication(job)}</td>
                                        <td>
                                            <Link to={`/jobs-monitored/${job.id}/alerts`} >
                                                <Button >
                                                    Alerts
                                                </Button>
                                            </Link>&nbsp;
                                            <Link to={`/jobs-monitored/${job.id}`} >
                                                <Button >
                                                    More details
                                                </Button>
                                            </Link>&nbsp;
                                            <CopyToClipboard 
                                            className="btn btn-primary"
                                            text={generateLinkJob(job)}
                                                >
                                                <button>Copy to link add your job</button>
                                            </CopyToClipboard>
                                        </td>
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