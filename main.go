package main

import (
	_ "encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	Database     = "geosn"
	SmCollection = "sm"
	GmCollection = "gm"
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
type UserLocations []UserLocation

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Coordinates struct {
	lat  float64
	long float64
}

func main() {
	cluster := "localhost"

	session, err := mgo.Dial(cluster)
	if err != nil {
		log.Fatal("could not connect to db: ", err)
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB(Database)
	coordinates := Coordinates{long: 3.575430130586027, lat: -65.18024627119303}
	scope := 50000000 // max distance in metres
	userid := 46
	k := 2

	ul := GetUserLocation(db, userid)
	gf := GetFriends(db, userid)
	notF := AreFriends(db, userid, 3) //false
	isF := AreFriends(db, userid, 45) //true

	rf := RangeFriends1(db, userid, coordinates, scope)
	nf := NearestFriends1(db, userid, coordinates, k)
	_ = RangeFriends2(db, userid, coordinates, scope)
	_ = NearestFriends2(db, userid, coordinates, k)
	_ = RangeFriends3(db, userid, coordinates, scope)
	_ = NearestFriends3(db, userid, coordinates, k)

	ru := RangeUsers(db, coordinates, scope)
	nu := NearestUsers(db, coordinates, k)

	fmt.Println("User with UserID:", userid)
	fmt.Println("Is at Coordinates:", ul.Location.Coordinates[0], ul.Location.Coordinates[1])
	fmt.Println("Has friends with UserIDs:", gf)
	fmt.Println("Is friends with UserID:3 =", notF)
	fmt.Println("Is friends with UserID:45 =", isF)
	fmt.Println("Has friends within", scope, "meters with UserIDs=", rf.GetUserIDs())
	fmt.Println("His", k, "-th nearest friend(s) have UserID(s)", nf.GetUserIDs())
	fmt.Println("")

	fmt.Println("Users within", scope, "meter are:")
	fmt.Println(len(ru))
	fmt.Println("")

	fmt.Println(k, "users nearest to", coordinates, "are:")
	fmt.Println(nu)
	fmt.Println("")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
