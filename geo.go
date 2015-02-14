package main

import (
	"fmt"
	"github.com/kellydunn/golang-geo"
	"gopkg.in/mgo.v2"
)

//User u, location q, radius r
//1. U = RangeUsers(q, r), R = ∅
//2. For each user ui ∈ U
//3.  If AreFriends(u, ui), add ui into R
//4. Return R
func RangeFriends(db *mgo.Database, userid int, coordinates Coordinates, r int) []int {
	var users []UserLocation
	rangeFriendsList := make([]int, 0, 1)

	users = RangeUsers(db, coordinates, r)

	for _, user := range users {
		if user.UserId == userid {
			rangeFriendsList = append(rangeFriendsList, user.UserId)
		}
	}

	return rangeFriendsList
}

/* Algorithm 3 (NF3) */
// 1. R = ∅
// 2. While |R| < k
// 3. ui = NextNearestUser(q)
// 4. If AreFriends(u, ui), add ui into R
// 5. Return R
func NearestFriends(db *mgo.Database, userid int, coordinates Coordinates, k int) ([]int, []UserLocation) {
	resultSet := make([]int, 0, 0)
	resultSetLocs := make([]UserLocation, 0, 0)
	nearestUserCount := 1

	for len(resultSet) < k {
		users := NearestUsers(db, coordinates, nearestUserCount)
		if nearestUserCount > len(users) {
			break
		}

		ui := users[nearestUserCount-1].UserId

		if AreFriends(db, userid, ui) {
			resultSet = append(resultSet, ui)
			resultSetLocs = append(resultSetLocs, users[nearestUserCount-1])
		}
		nearestUserCount++

	}

	return resultSet, resultSetLocs
}

// Input: Location q, positive integer m
// Output: Result set R
// 1. Initialize bs = ∞, bun = 0, R = ∅, Fseen = ∅, i = 1
// 2. u1 = NearestUsers(q, 1), NSGu1 = {u1} ∪ NF(q, u1, m − 1)
// 3. bs = adist(q, NSGu1)
// 4. While bun < bs
// 5.   ui = NextNearestUser(q)
// 6.   F = RF(ui, q, bs)
// 7.   NSGui = {ui} ∪ the m − 1 nearest users to q in F
// 8.   If |NSGui| = m ∧ adist(q, NSGui) < bs
// 9.     R = (NSGui), bs = adist(q, NSGui)
// 10.  For all u ∈ F
// 11.    If u 6∈ Fseen, add u to Fseen
// 12.    If |NSGu| < m, add ui to NSGu
// 13.    If |NSGu| = m ∧ adist(q, NSGu) < bs
// 14.      R = (NSGu), bs = adist(q, NSGu)
// 15.  bun = m · ||q, ui||
// 16.  i + +
// 17. i = i − 1 // so that ui is the lastly seen user
// 18. For all u ∈ Fseen ∧ |NSGu| < m
// 19.  If adist(q, NSGu) + (m − |NSGu|) · ||q, ui|| < bs
// 20.    NSGu = NF(u, q, m − 1)
// 21.    If adist(q, NSGu) < bs
// 22.      R = (NSGu), bs = adist(q, NSGu)
// 23. Return R
func NearestStarGroup(db *mgo.Database, q Coordinates, m int) {
	bs := 0.0 //Inf
	// bun := 0.0
	// resultSet := 0
	// fseen := 0
	// i := 1

	nu := NearestUsers(db, q, 1)
	u1 := nu[0].UserId

	NSGu1 := make([]UserLocation, 0, 1)
	_, nf := NearestFriends(db, u1, q, m-1)

	NSGu1 = append(NSGu1, nf...)

	for _, loc := range NSGu1 {
		dist := q.CalcDistance(loc) //find distance of each userloc to q point
		bs += dist                  //sum distances to bs
	}
	fmt.Println("Distance =", bs)
	// for bun < bs {

	// }

}

func (q Coordinates) CalcDistance(loc UserLocation) float64 {
	pointq := geo.NewPoint(q.long, q.lat)
	long := loc.Location.Coordinates[0]
	lat := loc.Location.Coordinates[1]
	pointLoc := geo.NewPoint(long, lat)
	dist := pointq.GreatCircleDistance(pointLoc)
	return dist
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
