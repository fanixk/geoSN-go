package main

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetFriends(c *mgo.Collection, userid int) []int {
	var result []User
	friends_list := make([]int, 0, 1)

	err := c.Find(bson.M{"userid": userid}).Select(bson.M{"friends_list": 1}).All(&result)
	if err != nil {
		panic(err)
	}

	for _, friend := range result {
		for _, f := range friend.Friends {
			friends_list = append(friends_list, f)
		}
	}
	return friends_list
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

func NearestUsers(c *mgo.Collection, long float64, lat float64, k int) []UserLocation {
	var res []UserLocation

	err := c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}).Limit(k).All(&res)

	if err != nil {
		panic(err)
	}

	return res
}

//User u, location q, radius r
//1.U = RangeUsers(q, r), R = ∅
//2. For each user ui ∈ U
//3. 	If AreFriends(u, ui), add ui into R
//4. Return R
func RangeFriends(c *mgo.Collection, userid int, long float64, lat float64, r int) []int {
	var users []UserLocation
	range_friends_list := make([]int, 0, 1)

	users = RangeUsers(c, long, lat, r)

	for _, user := range users {
		if user.UserId == userid {
			range_friends_list = append(range_friends_list, user.UserId)
		}
	}

	return range_friends_list
}
