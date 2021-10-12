import AbstractService from "./AbstractService";

export default class JobService extends AbstractService {
    
    getHistoryNotifications = async (jobId) => {
        return this.httpClient()
        .get(`${process.env.REACT_APP_URL_BASE}job-notifications-received/${jobId}`)
        .then(this.extractReponse)
    }

    getAllAlerts(jobId) {
        return this.httpClient()
            .get(
                `${process.env.REACT_APP_URL_BASE}companies/${this.getCompanyId()}/jobs/${jobId}/alerts`,
            )
            .then(this.extractReponse)
    }

    createAlert(alert) {
        return this.httpClient()
        .post(
            `${process.env.REACT_APP_URL_BASE}companies/${this.getCompanyId()}/jobs/${alert.jobId}/alerts`,
             alert
        )
        .then(this.extractReponse)
    }

    getAll() {
        return this.httpClient()
            .get(`${process.env.REACT_APP_URL_BASE}companies/${this.getCompanyId()}/jobs`)
            .then(this.extractReponse)
    }

    create(job) {
        return this.httpClient().post(
            `${process.env.REACT_APP_URL_BASE}companies/${this.getCompanyId()}/jobs`, job
        )
    }

    getEventNotifications() {
        return this.httpClient()
            .get(`${process.env.REACT_APP_URL_BASE}event-notifications/${this.getCompanyId()}`)
            .then(this.extractReponse)
    }

    generateLinkJob(job) {
        return (`${process.env.REACT_APP_URL_BASE}event-notifications/${this.getCompanyId()}/${job.id}`)
    }

}