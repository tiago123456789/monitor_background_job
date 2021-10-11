import React, { useEffect, useState } from "react"
import { Button, Container, Table } from "reactstrap"
import axios from "axios"
import moment from "moment"

function App(props) {
  const [jobHistories, setJobHistories] = useState([])

  const getHistoryNotifications = async (jobId) => {
    const registers = await axios.get(`http://localhost:4000/job-notifications-received/${jobId}`)
      .then(({ data }) => data)
    setJobHistories(registers)
  }

  useEffect(() => {
    (async () => {
      await getHistoryNotifications(props.match.params.id)
    })()

  }, [])

  return (
    <>
      <Container>
      {jobHistories.length == 0 && <p>No have registers</p>}
        {jobHistories.length > 0 &&
          (
            <>
              <h1>Job details</h1>
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
