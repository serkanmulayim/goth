import React from "react";
import AdminApi from "../../api/AdminApi";
import {Button, ButtonGroup} from "react-bootstrap";
import BootstrapTable from "react-bootstrap-table-next";
import BootstrapButton from "react-bootstrap-table-next";
import "react-bootstrap-table-next/dist/react-bootstrap-table2.min.css";
import AdminDeleteModal from "./AdminDeleteModal";
import AdminUpdateModal from "./AdminUpdateModal";
import AdminCreateModal from "./AdminCreateModal";

const tableCols = [
  {
    dataField: "email",
    text: "Name",
    sort: true,
    formatter:(cell, row, rowIndex, formatExtraData) => {
      return (<p>{row.firstname} {row.lastname}</p>);
    }
  },
  {
    dataField: "firstname",
    text: "email",
    sort: true,
    formatter:(cell, row, rowIndex, formatExtraData) => {
      return (<p><a href={"mailto:" + row.email}>{row.email}</a></p>);
    }
  },
  {
    dataField:"",
    text:"",
    formatter:(cellContent, row, rowIndex, formatExtraData) => {
      return (
        <ButtonGroup>
            <Button size="sm" type="button" variant="link" onClick={(e) => updateAdmin2(e, row)}>Update</Button>
            <Button size="sm" type="button" variant="link" onClick={(e) => deleteAdmin2(e, row)}>Delete</Button>
        </ButtonGroup>
      );
    }
  }
];
var deleteAdmin2, updateAdmin2;
function deleteAdmin(e, row) {
  e.preventDefault();
  this.setState({showCreateUpdateModal: false, showDeleteModal: true, affectedAdmin:row, showCreateModal:false});

}

function updateAdmin(e, row) {
  e.preventDefault();
  this.setState({showUpdateModal: true, showDeleteModal: false,affectedAdmin:row, showCreateModal:false});
}

class Admins extends React.Component{

  constructor(props) {
    super(props);
    deleteAdmin2 = deleteAdmin.bind(this);
    updateAdmin2 = updateAdmin.bind(this);
    this.state={admins:[], showDeleteModal:false, showCreateUpdateModal:false, showCreateModal:false, affectedAdmin:{}};
    this.messageText = React.createRef();
  }


  hideModals = (e, message) => {
    if(!!message) {
      console.log("message", message);
      this.messageText.current.innerText = message;
    } else {
      this.messageText.current.innerText = "";
    }
    this.setState({showDeleteModal:false, showUpdateModal:false, showCreateModal:false,  affectedAdmin:{}});
  }

  buttonFormatter = (cell, row) =>{
    return (<BootstrapButton type="submit">Update</BootstrapButton>);
  }

  createAdmin=(e) => {
    e.preventDefault();
    this.setState({showUpdateModal: false, showDeleteModal: false,showCreateModal:true,  affectedAdmin:{}});
  }

  render() {
    console.log("rerenderig", this.state);

    return (

      <div className="Admin" style={{width:"80%", margin:"20px auto", textAlign:"center", border:'1px solid gray'}}>
          <br/>
          <h2 className="h2">Admin Accounts</h2>
          <div >
              <BootstrapTable
                  keyField="email"
                  data={this.state.admins}
                  noDataIndication="No Admins? Weird..."
                  columns={tableCols}
                  selectRow={{
                      mode: "checkbox",
                      hideSelectColumn:true,
                      clickToSelect:false,
                      onSelect:this.selectRow
                  }}
                  hover
                  striped
                  condensed
                  bootstrap4
                  bordered
              />

          </div>

          <div style={{textAlign:"right", margin:"auto 5%"}}>
              <label style={{color:'green', margin:"auto 1%"}} ref={this.messageText}>{this.state.message}</label>
              <Button  variant="primary" onClick={this.createAdmin}>Create Account</Button>
              <AdminDeleteModal show={this.state.showDeleteModal} errorMessage="" onHideWrapper={this.hideModals} admin={this.state.affectedAdmin} />
              <AdminUpdateModal show={this.state.showUpdateModal} errorMessage="" onHideWrapper={this.hideModals} admin={this.state.affectedAdmin} />
              <AdminCreateModal show={this.state.showCreateModal} errorMessage="" onHideWrapper={this.hideModals}/>
          </div>

          <br/><br/>
      </div>
    );
  }

  componentDidMount() {
    AdminApi.GetAdmins(this.props.handleLogout)
    .then(resp => {
      this.setState({admins:resp.admins});

    });

  }

}

export default Admins;
