import React, { Component } from "react";
import L from "leaflet";

class MapComponent extends Component {
  componentDidMount() {
    const position = [51.505, -0.09];
    const map = L.map("map").setView(position, 13);

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution:
        '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    L.marker(position)
      .addTo(map)
      .bindPopup("A pretty CSS3 popup. <br> Easily customizable.");
  }

  render() {
    return (
      <div className="leaflet_popup">
        <div className="leaflet_popup_inner" id="map"></div>
      </div>
    );
  }
}

export default MapComponent;
