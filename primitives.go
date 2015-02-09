package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetFriends(c *mgo.Collection, userid int) []User {
	var result []User

	err := c.Find(bson.M{"userid": userid}).Select(bson.M{"friends_list": 1}).All(&result)
	if err != nil {
		panic(err)
	}
	return result
}

func AreFriends(c *mgo.Collection, userid1 int, userid2 int) bool {
	//we suppose if userid2 exists in users1 friends_list then the opposite holds true
	//db.sm.count({   "$and": [ {userid: userid1}, { "friends_list":  { "$in": [ userid2 ] }}] })
	count, err := c.Find(
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

func RangeUsers(c *mgo.Collection, long float64, lat float64, scope int) []UserLocation {
	var res []UserLocation

	err := c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
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
