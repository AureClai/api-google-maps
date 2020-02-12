import React, { Component } from "react";
import logo from "../img/Logo_Cerema.png";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";
import Socket from "../socket.js";
import SettingsComponent from "./Tabs/SettingsComponent.jsx";
import RequestTab from "./Tabs/RequestTab.jsx";
import ImportTab from "./Tabs/ImportTab";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      connected: true,
      settings: {
        api_key: "",
        proxy: ""
      },
      currentTab: "Request"
    };
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
  }

  onDisconnect() {
    this.setState({ connected: false });
    this.testLog("Websocket déconnecté... Relancez l'application...", "FAIL");
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
    const proxy = data["proxy"];
    const api_key = data["api-key"];
    this.setState({
      settings: {
        api_key: api_key,
        proxy: proxy
      }
    });
  }

  /*
    Tabs mangement
  */
  onClickParameters(e) {
    this.setState({ currentTab: "Parameters" });
  }

  onClickRequest(e) {
    this.setState({ currentTab: "Request" });
  }

  onClickImport(e) {
    this.setState({ currentTab: "Import" });
  }

  render() {
    var jsxTab;
    if (this.state.currentTab === "Parameters") {
      jsxTab = (
        <SettingsComponent
          {...this.state}
          onChangeSettings={this.onChangeSettings.bind(this)}
        />
      );
    } else if (this.state.currentTab === "Request") {
      jsxTab = (
        <RequestTab
          {...this.state}
          onTestRequest={this.onTestRequest.bind(this)}
        />
      );
    } else if (this.state.currentTab === "Import") {
      jsxTab = <ImportTab />;
    }

    return (
      <div className="App">
        <nav className="navbar navbar-light bg-light">
          <div id="cerema-brand">
            <img src={logo} className="App-logo" alt="logo" />
            <span className="navbar-brand h">Cerema API Google</span>
          </div>
        </nav>
        <div className="container-fluid fill-height">
          <div className="row flex-grow-1">
            <div id="tabExplorer" className="col-1">
              <li>
                <ul onClick={this.onClickParameters.bind(this)}>Paramètres</ul>
                <ul onClick={this.onClickImport.bind(this)}>Import</ul>
                <ul onClick={this.onClickRequest.bind(this)}>Requête</ul>
              </li>
            </div>
            {jsxTab}
          </div>
        </div>
        <div id="request-test" readOnly></div>
      </div>
    );
  }
}

export default App;
