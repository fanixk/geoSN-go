package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

/* create queries
 * use geosn
 * db.sm.insert({userid: 1, friends_list: [2] })
 * db.gm.insert({userid: 1,  "location" : { "type" : "Point", "coordinates" : [ 151.20699, -33.867487 ] } })
 * db.sm.ensureIndex({userid: true}, {unique: true})
 * db.gm.ensureIndex({userid: true}, {unique: true})
 * db.gm.ensureIndex({location:"2dsphere"})
 */

type User struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	UserId  int           `bson:"userid" json:"userid"`
	Friends []int         `bson:"friends_list" json:"friends_list"`
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

func main() {
	cluster := "localhost"

	session, err := mgo.Dial(cluster)
	if err != nil {
		log.Fatal("could not connect to db: ", err)
		panic(err)
	}

	defer session.Close()
	// session.SetMode(mgo.Monotonic, true)

	// query the database
	colgm := session.DB("geosn").C("gm")
	//colsm := session.DB("geosn").C("sm")

	long := 151.701642
	lat := -33.690647
	scope := 50000 // max distance in metres

	results := GetNearbyUsers(colgm, long, lat, scope)
	//results := GetFriends(colsm, 1)
	//results := AreFriends(colsm, 1, 3) //false
	//results := AreFriends(colsm, 1, 2) //true

	// convert it to JSON so it can be displayed
	formatter := json.MarshalIndent
	response, err := formatter(results, " ", "   ")

	fmt.Println(string(response))
}
