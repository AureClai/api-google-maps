package main

import (
	"encoding/json"
	"fmt"

	"googlemaps.github.io/maps"
)

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*Feature, 0),
	}
}

// AddFeature appends a feature to the collection.
func (fc *FeatureCollection) AddFeature(feature *Feature) *FeatureCollection {
	fc.Features = append(fc.Features, feature)
	return fc
}

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	type featureCollection FeatureCollection

	fcol := &featureCollection{
		Type: "FeatureCollection",
	}

	fcol.Features = fc.Features
	if fcol.Features == nil {
		fcol.Features = make([]*Feature, 0)
	}

	return json.Marshal(fcol)
}

type Feature struct {
	ID         interface{}            `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Geometry   *Geometry              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func NewFeature(geometry *Geometry) *Feature {
	return &Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(map[string]interface{}),
	}
}

// NewPointFeature creates and initializes a GeoJSON feature with a point geometry using the given coordinate.
func NewPointFeature(coordinate []float64) *Feature {
	return NewFeature(NewPointGeometry(coordinate))
}

// NewLineStringFeature creates and initializes a GeoJSON feature with a line string geometry using the given coordinates.
func NewLineStringFeature(coordinates [][]float64) *Feature {
	return NewFeature(NewLineStringGeometry(coordinates))
}

func (f Feature) MarshalJSON() ([]byte, error) {
	type feature Feature

	fea := &feature{
		ID:       f.ID,
		Type:     "Feature",
		Geometry: f.Geometry,
	}

	if f.Properties != nil && len(f.Properties) != 0 {
		fea.Properties = f.Properties
	}

	return json.Marshal(fea)
}

type GeometryType string

const (
	GeometryPoint      GeometryType = "Point"
	GeometryLineString GeometryType = "LineString"
)

type Geometry struct {
	Type       GeometryType `json:"type"`
	Point      []float64
	LineString [][]float64
}

func NewPointGeometry(coordinate []float64) *Geometry {
	return &Geometry{
		Type:  GeometryPoint,
		Point: coordinate,
	}
}

func NewLineStringGeometry(coordinates [][]float64) *Geometry {
	return &Geometry{
		Type:       GeometryLineString,
		LineString: coordinates,
	}
}

func (g Geometry) MarshalJSON() ([]byte, error) {
	// struct define the order of the JSON elements
	type geometry struct {
		Type        GeometryType `json:"type"`
		Coordinates interface{}  `json:"coordinates,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryLineString:
		geo.Coordinates = g.LineString
	}

	return json.Marshal(geo)
}

func convertLatLngToFloatArray(in []maps.LatLng) [][]float64 {
	output := make([][]float64, 0)
	for i := 0; i < len(in); i++ {
		output = append(output, []float64{in[i].Lng, in[i].Lat})
	}
	return output
}

func convertToGeoJSON(encodedPolyline string) string {
	decodedByMaps, _ := maps.DecodePolyline(encodedPolyline)
	// Creation of the json
	featureCollection := NewFeatureCollection()
	// Creation of the linestring feature
	featureCollection.AddFeature(NewLineStringFeature(convertLatLngToFloatArray(decodedByMaps)))
	outJson, err := json.Marshal(featureCollection)
	if err != nil {
		fmt.Println(err)
	}

	return string(outJson)
}
