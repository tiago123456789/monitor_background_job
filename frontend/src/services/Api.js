import axios from "axios"

axios.interceptors.request.use(function (config) {
  const token = localStorage.getItem("accessToken")
  config.headers.Authorization =  token;
  return config;
});

axios.interceptors.response.use(function (response) {
    return response;
  }, function (error) {
    if(error.response.status === 403) {  window.location.href = "/" }
    return Promise.reject(error);
});

export default axios;