import axios from 'axios'
import https from 'https'
//import {withRouter} from "react-router-dom";

const loginEndpoint = "/api/login";
const checkAuthEndpoint = "/api/checkauth";
const logoutEndpoint = "/api/logout";

// const loginEndpoint = "https://localhost.goth.com:8443/api/login";
// const checkAuthEndpoint = "https://localhost.goth.com:8443/api/checkauth";
// const logoutEndpoint = "https://localhost.goth.com:8443/api/logout";

// axios.defaults.withCredentials = true;
const axioss = axios.create({
  httpsAgent: new https.Agent({
      rejectUnauthorized: false
    }),
  withCredentials:true
});


const encodeForm = (data) => {
  return Object.keys(data)
      .map(key => encodeURIComponent(key) + '=' + encodeURIComponent(data[key]))
      .join('&');
}

class Transport {

  static async DoLoginRequest(username, password) {
    let form = encodeForm({username:username, password:password});
    try {
      let response = await axioss.post(loginEndpoint, form);
      localStorage.setItem("auth", true);
      return {success:true, user:response.data.user};

    } catch(error){
      return this.handleError(error);
    }
  }

  static async DoCheckAuthRequest(){
    try {
      let resp = await axioss.get(checkAuthEndpoint);
      return {success: true, status:200, user:resp.data.user};
    } catch (error) {
      localStorage.setItem("auth", true);
      return this.handleError(error);
    }
  }

  static async DoLogoutRequest(){
    try {
      await axioss.get(logoutEndpoint);
      localStorage.setItem("auth", false);
      return {success: true, status:200};
    } catch (error) {
      return this.handleError(error);
    }
  }



  static handleError(error) {

    if(!error.response ){
      return {success:false, status:-1};
    } else if(error.response.status === 401){
      return {success:false, status:error.response.status};
    } else {
      return {success:false, status:error.response.status};
    }
  }
}

export default Transport;
