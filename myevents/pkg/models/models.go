package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string        `json:"name"`
	Duration  int           `json:"duration"`
	StartDate int64         `json:"startdate"`
	EndDate   int64         `json:"enddate"`
	Location  Location      `json:"location"`
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string        `json:"name"`
	Address   string        `json:"address"`
	Country   string        `json:"country"`
	OpenTime  int           `json:"opentime"`
	CloseTime int           `json:"closetime"`
	Halls     []Hall        `json:"halls"`
}

type Hall struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `json:"name"`
	Location string        `json:"location,omitempty"`
	Capacity string        `json:"capacity"`
}
