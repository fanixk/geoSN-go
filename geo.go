package main

import (
	"github.com/kellydunn/golang-geo"
	"gopkg.in/mgo.v2"
	"time"
)

//User u, location q, radius r
//1. U = RangeUsers(q, r), R = ∅
//2. For each user ui ∈ U
//3.  If AreFriends(u, ui), add ui into R
//4. Return R
func RangeFriends(db *mgo.Database, userid int, coordinates Coordinates, r int) UserLocations {
	defer timeTrack(time.Now(), "RangeFriends")

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

/* Algorithm 3 (NF3) */
// 1. R = ∅
// 2. While |R| < k
// 3. ui = NextNearestUser(q)
// 4. If AreFriends(u, ui), add ui into R
// 5. Return R
func NearestFriends(db *mgo.Database, userid int, coordinates Coordinates, k int) UserLocations {
	// defer timeTrack(time.Now(), "NearestFriends")

	resultSet := make([]UserLocation, 0, 0)
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
	pointq := geo.NewPoint(q.long, q.lat)
	long := loc.Location.Coordinates[0]
	lat := loc.Location.Coordinates[1]
	pointLoc := geo.NewPoint(long, lat)
	dist := pointq.GreatCircleDistance(pointLoc)
	return dist
}
