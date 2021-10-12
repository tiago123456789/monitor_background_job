import React, { useEffect, useState } from "react"
import { Button, Container, Table } from "reactstrap"
import axios from "axios"
import moment from "moment"
import JobService from "../services/JobService"

const jobService = new JobService()

function App(props) {
  const [jobHistories, setJobHistories] = useState([])

  const getHistoryNotifications = async (jobId) => {
    const registers = await jobService.getHistoryNotifications(jobId);
    setJobHistories(registers)
  }

  useEffect(() => {
    (async () => {
      await getHistoryNotifications(props.match.params.id)
    })()

  }, [])

  const isEmpty = () => {
    return jobHistories.length == 0
  }

  return (
    <>
      <Container>
      { isEmpty() && <p>No have registers</p>}
      { !isEmpty() > 0 &&
          (
            <>
              <h1 className="text-center">Job details</h1>
              <Table>
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Occour At</th>
                  </tr>
                </thead>
                <tbody>
                  {
                    jobHistories.map(job => {
                      return (
                        <tr>
                          <th scope="row">{job.id}</th>
                          <td>{moment(job.OccourAt).format("DD/MM/YY HH:mm:ss")}</td>
                        </tr>
                      )
                    })
                  }
                </tbody>
              </Table>
            </>
          )
        }
      </Container>
    </>
  );
}

export default App;
