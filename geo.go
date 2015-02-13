package main

import (
	"gopkg.in/mgo.v2"
)

//User u, location q, radius r
//1. U = RangeUsers(q, r), R = ∅
//2. For each user ui ∈ U
//3.  If AreFriends(u, ui), add ui into R
//4. Return R
func RangeFriends(db *mgo.Database, userid int, long float64, lat float64, r int) []int {
	var users []UserLocation
	range_friends_list := make([]int, 0, 1)

	users = RangeUsers(db, long, lat, r)

	for _, user := range users {
		if user.UserId == userid {
			range_friends_list = append(range_friends_list, user.UserId)
		}
	}

	return range_friends_list
}
