import React from "react";
import AdminApi from "../../api/AdminApi";
import {Modal, Button} from "react-bootstrap" ;


class AdminDeleteModal extends React.Component {
  constructor(props) {
    super(props);
    this.state={};
    this.messageText= React.createRef();
  }

  handleDelete = (e) => {
    e.preventDefault();
    AdminApi.DeleteAdmin(null, this.props.admin)
    .then(resp => {
      if(resp.success){
        this.onHide(null, "Admin is deleted");
      } else {
        this.messageText.current.innerText= resp.message;
      }
    })
  }

  onHide = (e, messageToParent) => {
    this.messageText.current.innerText = "";
    this.props.onHideWrapper(null, messageToParent);
  }

  render() {
    if(!this.props.admin){
      return <div></div>
    }
    return(
      <div>
          <Modal size="lg" show={this.props.show} onHide={this.onHide}>
              <Modal.Header closeButton closeLabel="Close">
                  <Modal.Title>
                      Delete Admin
                  </Modal.Title>
              </Modal.Header>
              <Modal.Body>
                  Delete admin, <b>{this.props.admin.firstname} {this.props.admin.lastname}</b>, with email <a href={"mailto:" + this.props.admin.email}>{this.props.admin.email}</a>.
              </Modal.Body>
              <Modal.Footer>
                  <p style={{color:'red'}} ref={this.messageText}>{this.state.message}</p>
                  <Button variant="secondary" onClick={this.onHide}>Cancel</Button>
            <Button variant="primary" onClick={this.handleDelete}>Delete</Button>
          </Modal.Footer>
        </Modal>
      </div>

    )
  }
}

export default AdminDeleteModal;
