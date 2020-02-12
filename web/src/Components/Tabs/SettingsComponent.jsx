import React, { Component } from "react";
import PropTypes from "prop-types";

class SettingsComponent extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    var { api_key, proxy } = this.props.settings;
    var { connected } = this.props;

    return (
      <div id="tabContent" className="col-11">
        <form className="setting-form">
          <label htmlFor="api-key-input">Clé API</label>
          <br></br>
          <input
            type="text"
            id="api-key-input"
            placeholder="Ici, la clé API"
            defaultValue={api_key}
            disabled={!connected}
          ></input>
          <br></br>
          <label htmlFor="proxy-input">Proxy</label>
          <br></br>
          <input
            type="text"
            id="proxy-input"
            placeholder="Ici, le proxy"
            defaultValue={proxy}
            disabled={!connected}
          ></input>
          <br></br>
          <button
            onClick={this.props.onChangeSettings}
            className="App-button"
            id="settings-submit"
            disabled={!connected}
          >
            <span>Valider</span>
          </button>
        </form>
      </div>
    );
  }
}

Component.propTypes = {
  connected: PropTypes.bool,
  settings: PropTypes.object,
  onChangeSettings: PropTypes.func
};

export default SettingsComponent;
