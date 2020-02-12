import React, { Component } from "react";
import logo from "../img/Logo_Cerema.png";
import "./App.css";
import Socket from "../socket.js";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      connected: false,
      settings: {
        api_key: "",
        proxy: ""
      }
    };
    this.onChangeSettings = this.onChangeSettings.bind(this);
    this.onTestRequest = this.onTestRequest.bind(this);
  }

  /*
  TESTING SECTION !
  */

  onTestRequest(e) {
    var { connected } = this.state;
    if (connected) {
      this.testLog("Demande d'un test de requête", "INFO");
      this.socket.emit("test request", {});
    }
  }

  onTestCallback(data) {
    this.testLog(data["message"], data["type"]);
  }

  testLog(message, type) {
    const dateString = new Date().toLocaleString();
    document.getElementById("request-test").innerHTML =
      `<p class=${type}>[${dateString}] ${message}\n</p>` +
      document.getElementById("request-test").innerHTML;
  }

  /*

  */

  componentDidMount() {
    this.socket = new Socket(new WebSocket("ws://localhost:8080/ws"));
    let socket = this.socket;

    socket.on("connect", this.onConnect.bind(this));
    socket.on("disconnect", this.onDisconnect.bind(this));
    socket.on("settings change", this.onReceiveSettingsInfo.bind(this));
    socket.on("test callback", this.onTestCallback.bind(this));
  }

  onConnect() {
    this.setState({ connected: true });
    document.getElementById("api-key-input").disabled = false;
    document.getElementById("proxy-input").disabled = false;
    document.getElementById("settings-submit").disabled = false;
    document.getElementById("request-test-button").disabled = false;
  }

  onDisconnect() {
    this.setState({ connected: false });
    this.testLog("Deconnecté du Websocket... relancez l'appli", "FAIL");
    document.getElementById("api-key-input").disabled = true;
    document.getElementById("proxy-input").disabled = true;
    document.getElementById("settings-submit").disabled = true;
    document.getElementById("request-test-button").disabled = true;
    alert("Relancez l'application...");
  }

  onChangeSettings(e) {
    e.preventDefault();
    var { connected } = this.state;
    if (connected) {
      let newAPIkey = document.getElementById("api-key-input").value;
      let newProxy = document.getElementById("proxy-input").value;
      this.socket.emit("settings change", {
        "api-key": newAPIkey,
        proxy: newProxy
      });
      console.log("Send message !");
    }
  }

  onReceiveSettingsInfo(data) {
    document.getElementById("proxy-input").value = data["proxy"];
    document.getElementById("api-key-input").value = data["api-key"];
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <h1>Cerema API Google</h1>
          <img src={logo} className="App-logo" alt="logo" />
          <form className="setting-form">
            <label for="api-key-input">Clé API</label>
            <br></br>
            <input
              type="text"
              id="api-key-input"
              placeholder="Ici, la clé API"
              disabled
            ></input>
            <br></br>
            <label for="proxy-input">Proxy</label>
            <br></br>
            <input
              type="text"
              id="proxy-input"
              placeholder="Ici, le proxy"
              disabled
            ></input>
            <br></br>
            <button
              onClick={this.onChangeSettings}
              className="App-button"
              id="settings-submit"
            >
              <span>Valider</span>
            </button>
          </form>
          <button
            onClick={this.onTestRequest}
            className="App-button"
            id="request-test-button"
          >
            <span>Tester une requête</span>
          </button>
          <div id="request-test" readOnly></div>
        </header>
      </div>
    );
  }
}

export default App;
