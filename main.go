package main

import (
	_ "encoding/json"
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
	coordinates := Coordinates{long: 3.575430130586027, lat: -65.18024627119303}
	scope := 500 // max distance in metres
	userid := 46
	k := 2

	//TODO: return both users and list of user ids when needed

	ul := GetUserLocation(db, userid)
	gf := GetFriends(db, userid)
	not_f := AreFriends(db, userid, 3) //false
	is_f := AreFriends(db, userid, 45) //true
	rf := RangeFriends(db, userid, coordinates, scope)
	nf := NearestFriends(db, userid, coordinates, k)

	ru := RangeUsers(db, coordinates, scope)
	nu := NearestUsers(db, coordinates, k)

	//show actual users
	// users := GetUsers(db, results)

	// convert it to JSON so it can be displayed
	// formatter := json.MarshalIndent

	// response, _ := formatter(users, " ", "   ")
	// fmt.Println(string(response))

	fmt.Println("User with UserID:", userid)
	fmt.Println("Is at Coordinates:", ul.long, ul.lat)
	fmt.Println("Has friends with UserIDs:", gf)
	fmt.Println("Is friends with UserID:3 =", not_f)
	fmt.Println("Is friends with UserID:45 =", is_f)
	fmt.Println("Has friends within", scope, "meters with UserIDs=", rf)
	fmt.Println("His", k, "-th nearest friend(s) have UserID(s)", nf)
	fmt.Println("")

	fmt.Println("Users within", scope, "meter are:")
	fmt.Println(ru)
	fmt.Println("")

	fmt.Println(k, "users nearest to", coordinates, "are:")
	fmt.Println(nu)
	fmt.Println("")
}
