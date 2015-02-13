package main

import (
	// "fmt"
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

/* Algorithm 3 (NF3) */
// 1. R = ∅
// 2. While |R| < k
// 3. ui = NextNearestUser(q)
// 4. If AreFriends(u, ui), add ui into R
// 5. Return R
func NearestFriends(db *mgo.Database, userid int, long float64, lat float64, k int) []int {
	result_set := make([]int, 0, 1)
	users := NearestUsers(db, long, lat, k)
	index := 0

	for len(result_set) < k {
		if index == len(users) {
			break
		}

		ui := users[index].UserId
		index++

		if AreFriends(db, userid, ui) {
			result_set = append(result_set, ui)
		}
	}

	return result_set
}
