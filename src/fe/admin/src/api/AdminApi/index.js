import axios from 'axios'
import https from 'https'
//import {withRouter} from "react-router-dom";


const isDev= (process.env.NODE_ENV === "development");
const server = isDev ? "https://localhost.goth.com:8444" : "";
const loginPath = "/api/login";
const checkauthPath = "/api/checkauth";
const logoutPath = "/api/logout";
const adminsPath = "/api/admins";

const loginEndpoint = server + loginPath;
const checkAuthEndpoint = server + checkauthPath;
const logoutEndpoint = server + logoutPath;
const adminsEndpoint = server + adminsPath;

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

class AdminApi {

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

  static async GetAdmins(logoutFunction){
    try {
      let resp = await axioss.get(adminsEndpoint);
      return {success: true, status:200, admins:resp.data.admins};
    } catch (error) {
      if(!!logoutFunction) {
        logoutFunction();
      }
      localStorage.removeItem("auth");
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

export default AdminApi;
