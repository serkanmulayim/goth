import React from "react";
import AdminApi from "../../api/AdminApi";
import {Modal, Button, Form} from "react-bootstrap" ;


class AdminCreateModal extends React.Component {
  constructor(props) {
    super(props);
    this.state= {};
    this.labelMessage = React.createRef();
  }

  validateFields(event) {
    if (event.target.password.value !== event.target.verifypassword.value) {
      return "Passwords do not match";
    }
    return null;
}

  handleCreate = (e) => {
    e.preventDefault();

    let err = this.validateFields(e);
    if(err != null) {
      this.labelMessage.current.innerText= err;
      return;
    }
    AdminApi.CreateAdmin(null, this.props.admin)
    .then(resp => {
        console.log("resp" , resp);
        if(resp.success){
          this.onHide(null, "Admin is created");
        } else {
          this.labelMessage.current.innerText= resp.message;
        }
      });
  }

  onHide = (e, messageToParent) => {
    this.labelMessage.current.innerText= "";
    this.props.onHideWrapper(null, messageToParent);
  }

  render() {
    return(
        <div>
            <Modal size="lg" show={this.props.show} onHide={this.onHide} >
                <Form onSubmit={this.handleCreate} autoComplete="false">
                    <Modal.Header closeButton closeLabel="Close">
                        <Modal.Title>Create Admin</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {/* {this.getFormGroup({as:Form.Col, controlId:"firstname", size:"large"}, "FirstName", {autoFocus:true, required:true, type:"input", placeHolder:"FirstName", defaultValue:(!!this.props.admin)? this.props.admin.firstname : ""})} */}
                        <Form.Group as={Form.Col} controlId="firstname" size="large">
                            <Form.Label>Firstname</Form.Label>
                            <Form.Control  autoFocus required type="input" placeholder="Firstname" />
                        </Form.Group>
                        <Form.Group as={Form.Col} controlId="lastname" size="large">
                            <Form.Label>Lastname</Form.Label>
                            <Form.Control  required type="input" placeholder="Lastname"/>
                        </Form.Group>
                        <Form.Group  controlId="email" size="large">
                            <Form.Label>email</Form.Label>
                            <Form.Control required type="text" placeholder="email" />
                        </Form.Group>
                        <Form.Group  controlId="phone" size="large">
                            <Form.Label>Phone</Form.Label>
                            <Form.Control autoFocus type="input" placeholder="555-555-5555" />
                        </Form.Group>
                        <Form.Group autoComplete="false"  controlId="address" size="large">
                            <Form.Label>Address</Form.Label>
                            <Form.Control autoComplete="false" type="text" placeholder="Address" />
                        </Form.Group>


                        <Form.Group controlId="password" size="large" >
                            <Form.Label>Password</Form.Label>
                            <Form.Control autoComplete="new-password" required type="password" placeholder="Password"/>
                        </Form.Group>

                        <Form.Group  controlId="verifypassword" size="large">
                            <Form.Label>Verify Password</Form.Label>
                            <Form.Control autoComplete="new-password" required type="password" placeholder="Password"/>
                        </Form.Group>

                    </Modal.Body>
                    <Modal.Footer>
                        <Form.Label id="errorLabel" style={{color:'red'}} ref={this.labelMessage}>{this.state.message}</Form.Label>
                        <Button variant="primary" type="submit">Create</Button>
                        <Button variant="secondary" onClick={this.onHide}>Cancel</Button>
                    </Modal.Footer>
                </Form>
            </Modal>
        </div>

    )
  }
}

export default AdminCreateModal;
