import React from "react";
import AdminApi from "../../api/AdminApi";
import {Button, ButtonToolbar, ButtonGroup} from "react-bootstrap";
//import {TableHeaderColumn, BootstrapButton} from "react-bootstrap-table-next";
import BootstrapTable from "react-bootstrap-table-next";
import TableHeaderColumn from "react-bootstrap-table-next";
import BootstrapButton from "react-bootstrap-table-next";
import "react-bootstrap-table-next/dist/react-bootstrap-table2.min.css";


class Admins extends React.Component{
tableCols = [
  {
    dataField: "firstname",
    text: "Firstname",
    sort: true
  },
  {
    dataField: "lastname",
    text: "Lastname",
    sort: true
  },
  {
    dataField: "email",
    text: "Email",
    sort: true
  },
  {
    dataField:"",
    text:"",
    formatter:(cellContent, row, rowIndex, formatExtraData) => {
      //return <a onClick={alert("clicked")}>row.firstname</a>
      //alert(JSON.stringify(row))
       return (

        <ButtonGroup vertical>
          <Button size="sm" type="button" variant="link" onClick={this.updateUser(row)}>Update</Button>
          <Button size="sm" type="button" variant="link" onClick={this.deleteUser(row)}>Delete</Button>
        </ButtonGroup>

      );
    }
  }
];
  constructor(props) {
    super(props);
    console.log(props);
    this.state={admins:[]};
  }

  updateUser= (row) => {
    return function(e) {
      alert(row.firstname + " update")
    }
  }

  deleteUser=(row) =>{
    let x = this.state.admins.length
    return function(e) {
      alert(row.firstname + " delete "+ x)
    }
  }

  selectRow = (row, isSelected, rowIndex) =>{
    alert(JSON.stringify(row));
  }

  buttonFormatter = (cell, row) =>{
    return (<BootstrapButton type="submit">Update</BootstrapButton>);
  }

  render() {

    return (

      <div style={{width:"80%", margin:"20px auto", textAlign:"center", border:'1px solid gray'}}>
          <br/>
          <h2 className="h2">Admin Accounts</h2>
          <div >

          <BootstrapTable
              keyField="email"
              data={this.state.admins}
              noDataIndication="No Admins? Weird..."
              columns={this.tableCols}
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
              bordered/>

          </div>

          <div style={{textAlign:"right", margin:"auto 5%"}}>
              <Button  size="small" variant="secondary">Create Account</Button>
          </div>

          <br/><br/>
      </div>
    );
  }

  componentDidMount() {
    AdminApi.GetAdmins(this.props.handleLogout)
    .then(resp => {
      let newState = Object.assign({}, this.state);
      newState.admins = resp.admins;
      this.setState(newState);

    });

  }

}

export default Admins;
