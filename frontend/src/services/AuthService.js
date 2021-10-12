import AbstractService from "./AbstractService";

export default class AuthService extends AbstractService {
    
    authenticate(credential) {
        return this.httpClient()
            .post(`${process.env.REACT_APP_URL_BASE}auth/login`, credential)
            .then(this.extractReponse)
    }

    create(newRegister) {
        return this.httpClient()
            .post(`${process.env.REACT_APP_URL_BASE}companies`, newRegister)
    }

    isAuthenticated() {
        return localStorage.getItem("accessToken") != null
    }
}