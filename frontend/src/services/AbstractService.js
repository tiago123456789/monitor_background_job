import token from "../utils/token";
import api from "./Api"

export default class AbstractService {

    httpClient() {
        return api
    }

    extractReponse(response) {
        if (response.data) {
            return response.data;
        }
        return response;
    }

    getCompanyId() {
        const accessToken = localStorage.getItem("accessToken")
        return token.getValueInPayloud("companyId", accessToken)
    }
}