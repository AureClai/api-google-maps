import React, { Component } from "react";
import PropTypes from "prop-types";
import {
  Button,
  Modal,
  Form,
  Row,
  OverlayTrigger,
  Tooltip,
  Table
} from "react-bootstrap";
import MapComponent from "./Popups/MapComponent.jsx";

class ImportTab extends Component {
  constructor(props) {
    super(props);
    this.state = {
      show: false,
      origLat: "45.6995715405897",
      origLong: "4.8357263290446",
      destLat: "45.5909623228595",
      destLong: "4.797397619176951",
      name: "test1"
    };
  }

  handleClose() {
    this.setState({
      show: false
    });
  }

  handleValidate() {
    this.props.addNewPath(this.state);
    this.handleClose.bind(this)();
  }

  handleMap(e) {
    this.props.askForPathMap(e.target.name.slice(5));
  }

  handleRemove(e) {
    this.props.askForRemovingPath(e.target.name.slice(8));
  }

  handleShow() {
    this.setState({
      show: true
    });
  }

  handleInputsChange(event) {
    var name = event.target.name;
    var value = event.target.value;
    this.setState({
      [name]: value
    });
  }

  togglePopup() {
    this.setState({
      showPopup: !this.state.showPopup
    });
  }

  render() {
    const renderTooltip = props => {
      return (
        <Tooltip {...props}>
          Ajouter un <b>itinéraire</b>
        </Tooltip>
      );
    };

    const renderTable = () => {
      var { paths } = this.props;
      return (
        <Table striped bordered hover size="sm">
          <thead>
            <tr>
              <th>#</th>
              <th>Nom</th>
              <th>Lat. Orig.</th>
              <th>Long. Orig.</th>
              <th>Lat. Dest.</th>
              <th>Long. Dest.</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {paths.map(path => {
              return (
                <tr>
                  <td>{path.id}</td>
                  <td>{path.name}</td>
                  <td>{path.origLat}</td>
                  <td>{path.origLong}</td>
                  <td>{path.destLat}</td>
                  <td>{path.destLong}</td>
                  <td>
                    <Button
                      variant="info"
                      name={"bMap_" + path.id}
                      onClick={this.handleMap.bind(this)}
                    >
                      Map
                    </Button>
                    <Button
                      variant="danger"
                      name={"bRemove_" + path.id}
                      onClick={this.handleRemove.bind(this)}
                    >
                      X
                    </Button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </Table>
      );
    };

    return (
      <div id="tabContent" className="col-11">
        <OverlayTrigger
          placement="right"
          delay={{ show: 250, hide: 50 }}
          overlay={renderTooltip}
        >
          <Button variant="primary" onClick={this.handleShow.bind(this)}>
            Ajouter
          </Button>
        </OverlayTrigger>
        {/* Here the table management */}
        <div id="pathTable">{renderTable()}</div>
        {this.state.show && <MapComponent></MapComponent>}
        {/*
        <Modal
          show={this.state.show}
          onHide={this.handleClose.bind(this)}
          aria-labelledby="contained-modal-title-vcenter"
          centered
        >
          <Modal.Header closeButton>
            <Modal.Title>Ajouter un itinéraire</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Form>
              <Form.Group as={Row}>
                <Form.Label>
                  <b>Nom :</b>
                </Form.Label>
                <Form.Control
                  type="input"
                  placeholder="Noms"
                  value={this.state.name}
                  name="name"
                  onChange={this.handleInputsChange.bind(this)}
                ></Form.Control>
                <Form.Label>
                  <b>Origine :</b>
                </Form.Label>
                <Form.Control
                  type="input"
                  placeholder="Latitude"
                  value={this.state.origLat}
                  name="origLat"
                  onChange={this.handleInputsChange.bind(this)}
                ></Form.Control>
                <Form.Control
                  type="input"
                  placeholder="Longitude"
                  value={this.state.origLong}
                  name="origLong"
                  onChange={this.handleInputsChange.bind(this)}
                ></Form.Control>
              </Form.Group>
              <Form.Group as={Row}>
                <Form.Label>
                  <b>Destination :</b>
                </Form.Label>
                <Form.Control
                  type="input"
                  placeholder="Latitude"
                  value={this.state.destLat}
                  name="destLat"
                  onChange={this.handleInputsChange.bind(this)}
                ></Form.Control>
                <Form.Control
                  type="input"
                  placeholder="Longitude"
                  value={this.state.destLong}
                  name="destLong"
                  onChange={this.handleInputsChange.bind(this)}
                ></Form.Control>
              </Form.Group>
            </Form>
          </Modal.Body>
          <Modal.Footer>
            <Button
              className="mx-2"
              variant="secondary"
              onClick={this.handleClose.bind(this)}
            >
              Fermer
            </Button>
            <Button variant="primary" onClick={this.handleValidate.bind(this)}>
              Sauvegarder les changements
            </Button>
          </Modal.Footer>
        </Modal>*/}
      </div>
    );
  }
}

Component.propTypes = {};

export default ImportTab;
