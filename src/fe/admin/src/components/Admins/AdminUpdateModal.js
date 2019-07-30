import React from "react";
import AdminApi from "../../api/AdminApi";
import {Modal, Button, Form} from "react-bootstrap" ;


class AdminUpdateModal extends React.Component {
  constructor(props) {
    super(props);
    this.state= {};
    this.labelMessage = React.createRef();
  }

  handleUpdate = (e) => {
    e.preventDefault();

    AdminApi.PutAdmin(null, this.props.admin)
    .then(resp => {
        console.log("resp", resp);
        if(resp.success){
          this.onHide(null, "Admin is updated");
        } else {
          this.labelMessage.current.innerText= resp.message;
        }
    });

  }

  onHide = (e, messageToParent) => {
    this.labelMessage.current.innerText= "";
    this.props.onHideWrapper(null, messageToParent);
  }

  getTitle= () => {
    return "Update : \"" + this.props.admin.firstname + " " + this.props.admin.lastname + "\"";
  }

  render() {
    return(
        <div>
            <Modal size="lg" show={this.props.show} onHide={this.onHide} >
                <Form onSubmit={this.handleUpdate} autoComplete="false">
                    <Modal.Header closeButton closeLabel="Close">
                        <Modal.Title>{this.getTitle()}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {/* {this.getFormGroup({as:Form.Col, controlId:"firstname", size:"large"}, "FirstName", {autoFocus:true, required:true, type:"input", placeHolder:"FirstName", defaultValue:(!!this.props.admin)? this.props.admin.firstname : ""})} */}
                        <Form.Group as={Form.Col} controlId="firstname" size="large">
                            <Form.Label>Firstname</Form.Label>
                            <Form.Control  autoFocus required type="input" placeholder="Firstname" defaultValue={this.props.admin.firstname}/>
                        </Form.Group>
                        <Form.Group as={Form.Col} controlId="lastname" size="large">
                            <Form.Label>Lastname</Form.Label>
                            <Form.Control  required type="input" placeholder="Lastname" defaultValue={this.props.admin.lastname}/>
                        </Form.Group>
                        <Form.Group  controlId="email" size="large">
                            <Form.Label>email</Form.Label>
                            <Form.Control required type="text" placeholder="email" defaultValue={this.props.admin.email}/>
                        </Form.Group>
                        <Form.Group  controlId="phone" size="large">
                            <Form.Label>Phone</Form.Label>
                            <Form.Control autoFocus type="input" placeholder="555-555-5555" defaultValue={this.props.admin.phone}/>
                        </Form.Group>
                        <Form.Group autoComplete="false"  controlId="address" size="large">
                            <Form.Label>Address</Form.Label>
                            <Form.Control autoComplete="false" type="text" placeholder="Address" defaultValue={this.props.admin.address}/>
                        </Form.Group>

                    </Modal.Body>
                    <Modal.Footer>
                        <Form.Label id="errorLabel" style={{color:'red'}} ref={this.labelMessage}>{this.state.message}</Form.Label>
                        <Button variant="primary" type="submit" >Update</Button>
                        <Button variant="secondary" onClick={this.onHide}>Cancel</Button>
                    </Modal.Footer>
                </Form>
            </Modal>
        </div>

    )
  }
}

export default AdminUpdateModal;
