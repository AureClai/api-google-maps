import React, { Component } from "react";
import PropTypes from "prop-types";

class RequestTab extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    var { api_key, proxy } = this.props.settings;
    var { connected } = this.props;

    return (
      <div id="tabContent" className="col-11">
        <button
          onClick={this.props.onTestRequest}
          className="App-button"
          id="request-test-button"
          disabled={!connected}
        >
          <span>Tester une requête</span>
        </button>
      </div>
    );
  }
}

Component.propTypes = {
  connected: PropTypes.bool,
  onTestRequest: PropTypes.func
};

export default RequestTab;
