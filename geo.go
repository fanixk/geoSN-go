package main

import (
	"github.com/kellydunn/golang-geo"
	"gopkg.in/mgo.v2"
	"sort"
	"time"
)

/* Algorithm 1 (RF1) */
// User u, location q, radius r
// 1. F = GetFriends(u), R = ∅
// 2. For each user ui ∈ F
// 3. GetUserLocation(ui)
// 4. If ||q, ui|| ≤ r, add ui into R
// 5. Return R
func RangeFriends1(db *mgo.Database, userid int, coordinates Coordinates, r int) UserLocations {
	defer timeTrack(time.Now(), "RangeFriends1")

	friends := GetFriends(db, userid)
	rangeFriendsList := make(UserLocations, 0, 1)

	for _, friend := range friends {
		userLocation := GetUserLocation(db, friend)
		dist := coordinates.CalcDistance(userLocation)
		if dist <= float64(r) {
			rangeFriendsList = append(rangeFriendsList, userLocation)
		}
	}
	return rangeFriendsList
}

/* Algorithm 2 (RF2) */
// 1. Return R = GetFriends(u) intersect RangeUsers(q, r)
func RangeFriends2(db *mgo.Database, userid int, coordinates Coordinates, r int) UserLocations {
	defer timeTrack(time.Now(), "RangeFriends2")

	rangeFriendsList := make(UserLocations, 0, 1)
	friends := GetFriends(db, userid)
	rangeUsers := RangeUsers(db, coordinates, r)

	for _, friend := range friends {
		for _, user := range rangeUsers {
			if friend == user.UserId {
				rangeFriendsList = append(rangeFriendsList, user)
			}
		}
	}

	return rangeFriendsList
}

// Algorithm 3 (RF3)
// User u, location q, radius r
// 1. U = RangeUsers(q, r), R = ∅
// 2. For each user ui ∈ U
// 3.  If AreFriends(u, ui), add ui into R
// 4. Return R
func RangeFriends3(db *mgo.Database, userid int, coordinates Coordinates, r int) UserLocations {
	defer timeTrack(time.Now(), "RangeFriends3")

	var users []UserLocation
	rangeFriendsList := make(UserLocations, 0, 1)

	users = RangeUsers(db, coordinates, r)

	for _, user := range users {
		if AreFriends(db, userid, user.UserId) {
			rangeFriendsList = append(rangeFriendsList, user)
		}
	}

	return rangeFriendsList
}

/* Algorithm 1 (NF1) */
// 1. F = GetFriends(u), R = ∅
// 2. For each user ui ∈ F, compute GetUserLocation(ui)
// 3. Sort F in ascending order of ||q, ui||
// 4. Insert the first k entries of F into R
// 5. Return R
func NearestFriends1(db *mgo.Database, userid int, coordinates Coordinates, k int) UserLocations {
	defer timeTrack(time.Now(), "NearestFriends1")

	resultSet := make([]UserLocation, 0, 1)
	nf1Map := make(map[float64]UserLocation)
	friends := GetFriends(db, userid)

	for _, friend := range friends {
		userLocation := GetUserLocation(db, friend)
		dist := coordinates.CalcDistance(userLocation)
		nf1Map[dist] = userLocation
	}

	var keys []float64
	for key := range nf1Map {
		keys = append(keys, key)
	}
	sort.Float64s(keys)

	for i := 0; i < k; i++ {
		resultSet = append(resultSet, nf1Map[keys[i]])
	}

	return resultSet
}

/* Algorithm 2 (NF2) */
// 1. F = GetFriends(u), R = ∅
// 2. While |R| < k
// 3. ui = NextNearestUser(q)
// 4. If ui ∈ F, add ui into R
// 5. Return R
func NearestFriends2(db *mgo.Database, userid int, coordinates Coordinates, k int) UserLocations {
	defer timeTrack(time.Now(), "NearestFriends2")

	friends := GetFriends(db, userid)
	resultSet := make([]UserLocation, 0, 1)
	nearestUserCount := 1

	for len(resultSet) < k {
		users := NearestUsers(db, coordinates, nearestUserCount)
		if nearestUserCount > len(users) {
			break
		}

		user := users[nearestUserCount-1]

		for _, friend := range friends {
			if friend == user.UserId {
				resultSet = append(resultSet, user)
			}
		}

		nearestUserCount++
	}
	return resultSet
}

/* Algorithm 3 (NF3) */
// 1. R = ∅
// 2. While |R| < k
// 3. ui = NextNearestUser(q)
// 4. If AreFriends(u, ui), add ui into R
// 5. Return R
func NearestFriends3(db *mgo.Database, userid int, coordinates Coordinates, k int) UserLocations {
	defer timeTrack(time.Now(), "NearestFriends3")

	resultSet := make([]UserLocation, 0, 1)
	nearestUserCount := 1

	for len(resultSet) < k {
		users := NearestUsers(db, coordinates, nearestUserCount)
		if nearestUserCount > len(users) {
			break
		}

		ui := users[nearestUserCount-1].UserId

		if AreFriends(db, userid, ui) {
			resultSet = append(resultSet, users[nearestUserCount-1])
		}
		nearestUserCount++

	}

	return resultSet
}

func (q Coordinates) CalcDistance(loc UserLocation) float64 {
	pointq := geo.NewPoint(q.lat, q.long)
	lat := loc.Location.Coordinates[0]
	long := loc.Location.Coordinates[1]
	pointLoc := geo.NewPoint(lat, long)
	dist := pointq.GreatCircleDistance(pointLoc)
	return dist
}

// func (q Coordinates) CalcDistance(loc Coordinates) float64 {
// 	pointq := geo.NewPoint(q.lat, q.long)
// 	pointLoc := geo.NewPoint(loc.lat, loc.long)
// 	dist := pointq.GreatCircleDistance(pointLoc)
// 	return dist
// }
