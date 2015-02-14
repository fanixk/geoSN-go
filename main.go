package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	DATABASE      = "geosn"
	SM_COLLECTION = "sm"
	GM_COLLECTION = "gm"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	UserId   int           `bson:"userid" json:"userid"`
	UserName string        `bson:"username" json:"username"`
	Friends  []int         `bson:"friends_list" json:"friends_list"`
}

type UserLocation struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	UserId   int           `bson:"userid" json:"userid"`
	Location GeoJson       `bson:"location" json:"location"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Coordinates struct {
	long float64
	lat  float64
}

func main() {
	cluster := "localhost"

	session, err := mgo.Dial(cluster)
	if err != nil {
		log.Fatal("could not connect to db: ", err)
		panic(err)
	}

	defer session.Close()
	// session.SetMode(mgo.Monotonic, true)

	db := session.DB(DATABASE)
	coordinates := Coordinates{long: 151.701642, lat: -33.690647}
	_ = coordinates
	scope := 50000 // max distance in metres
	_ = scope
	userid := 10
	_ = userid

	//TODO:
	//implement method that returns users when given an array of userids
	//to be used as utility method

	// results := RangeFriends(db, userid, coordinates, scope)
	// results := GetFriends(db, 1)
	// results := GetUserLocation(db, 2)
	// results := AreFriends(db, 1, 3) //false
	// results := AreFriends(colsm, 1, 2) //true
	// results := RangeUsers(db, coordinates, scope)
	// results := NearestUsers(db, coordinates, 3)
	results := NearestFriends(db, 45, coordinates, 3)

	// convert it to JSON so it can be displayed
	formatter := json.MarshalIndent
	response, err := formatter(results, " ", "   ")

	fmt.Println(string(response))
}
