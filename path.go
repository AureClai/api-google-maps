package main

import (
	"github.com/mitchellh/mapstructure"
)

type PathInfos struct {
	Coordinates struct {
		Origin struct {
			Lat  string
			Long string
		}
		Destination struct {
			Lat  string
			Long string
		}
	}
	Name string
	ID   int
}

func addNewPathToDatabase(infos interface{}) {
	newPathInfos := new(PathInfos)
	mapstructure.Decode(infos, &newPathInfos)
	go makeDirectionRequest(newPathInfos)
}
