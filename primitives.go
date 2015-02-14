package main

import (
	_ "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetFriends(db *mgo.Database, userid int) []int {
	var result []User
	collection := db.C(SmCollection)
	friendsList := make([]int, 0, 1)

	err := collection.Find(bson.M{"userid": userid}).Select(bson.M{"friends_list": 1}).All(&result)
	if err != nil {
		panic(err)
	}

	for _, friend := range result {
		for _, f := range friend.Friends {
			friendsList = append(friendsList, f)
		}
	}
	return friendsList
}

func AreFriends(db *mgo.Database, userid1 int, userid2 int) bool {
	collection := db.C(SmCollection)
	//we suppose if userid2 exists in users1 friends_list then the opposite holds true
	//db.sm.count({   "$and": [ {userid: userid1}, { "friends_list":  { "$in": [ userid2 ] }}] })
	count, err := collection.Find(
		bson.M{
			"$and": []bson.M{
				bson.M{
					"userid": userid1,
				},
				bson.M{
					"friends_list": bson.M{
						"$in": []int{userid2},
					},
				},
			},
		}).Count()

	if err != nil {
		panic(err)
	}

	if count > 0 {
		return true
	}
	return false
}

func RangeUsers(db *mgo.Database, coordinates Coordinates, scope int) []UserLocation {
	var res []UserLocation
	collection := db.C(GmCollection)

	err := collection.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{coordinates.long, coordinates.lat},
				},
				"$maxDistance": scope,
			},
		},
	}).All(&res)

	if err != nil {
		panic(err)
	}

	return res
}

func NearestUsers(db *mgo.Database, coordinates Coordinates, k int) []UserLocation {
	var res []UserLocation
	collection := db.C(GmCollection)

	err := collection.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{coordinates.long, coordinates.lat},
				},
			},
		},
	}).Limit(k).All(&res)

	if err != nil {
		panic(err)
	}

	return res
}

func GetUsers(db *mgo.Database, userids []int) []User {
	var users []User
	collection := db.C(GmCollection)
	err := collection.Find(
		bson.M{
			"userid": bson.M{
				"$in": userids,
			},
		},
	).All(&users)

	if err != nil {
		panic(err)
	}

	return users
}

func GetUserLocation(db *mgo.Database, userid int) Coordinates {
	collection := db.C(GmCollection)
	var location UserLocation
	err := collection.Find(bson.M{"userid": userid}).One(&location)

	if err != nil {
		panic(err)
	}

	coordinates := Coordinates{
		long: location.Location.Coordinates[0],
		lat:  location.Location.Coordinates[1],
	}
	return coordinates
}

func (ul UserLocations) GetUserIDs() []int {
	userIDs := make([]int, 0, 0)
	for _, user := range ul {
		userIDs = append(userIDs, user.UserId)
	}
	return userIDs
}
