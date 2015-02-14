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
func RangeFriends(db *mgo.Database, userid int, coordinates Coordinates, r int) []int {
	var users []UserLocation
	range_friends_list := make([]int, 0, 1)

	users = RangeUsers(db, coordinates, r)

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
func NearestFriends(db *mgo.Database, userid int, coordinates Coordinates, k int) []int {
	result_set := make([]int, 0, 0)
	nearest_user_count := 1

	for len(result_set) < k {
		users := NearestUsers(db, coordinates, nearest_user_count)
		if nearest_user_count > len(users) {
			break
		}

		ui := users[nearest_user_count-1].UserId
		nearest_user_count++

		if AreFriends(db, userid, ui) {
			result_set = append(result_set, ui)
		}
	}

	return result_set
}

// Algorithm 1 (NF1)
// 1. F = GetFriends(u), R = ∅
// 2. For each user ui ∈ F, compute GetUserLocation(ui)
// 3. Sort F in ascending order of ||q, ui||
// 4. Insert the first k entries of F into R
// 5. Return R
// func NearestFriends(db *mgo.Database, userid int, coordinates Coordinates, k int) []int {
//  friends := GetFriends(db, userid)
//  friends_locations := make([]UserLocation, 0, 0)

//  for _, friend := range friends {
//    user_location := GetUserLocation(db, friend)
//    friends_locations = append(friends_locations, user_location)
//  }

//  //TODO: sort

//  for i := 0; i < k; i++ {
//    result_set = append(result_set, friends_locations[i])
//  }

//  return result_set
// }
