use geosn;

db.sm.drop();
db.gm.drop();

for (var i = 1; i <= 20000; i++) {
  var friends_list = [];
  switch (i) {
    case 1:
      friends_list = [i + 1, 20000];
      break;
    case 20000:
      friends_list = [1, i - 1];
      break;
    default:
      friends_list = [i - 1, i + 1];
  }

  db.sm.insert({
    "userid": i,
    "username": "User" + i,
    "friends_list": friends_list
  });
}

for (i = 1; i <= 20000; i++) {
  var random_long = Math.random() * (180 + 180) - 180;
  var random_lat = Math.random() * (90 + 90) - 90;

  db.gm.insert({
    "userid": i,
    "location": {
      "type": "Point",
      "coordinates": [random_long, random_lat]
    }
  });
}

db.sm.ensureIndex({
  userid: true
}, {
  unique: true
});
db.gm.ensureIndex({
  userid: true
}, {
  unique: true
});
db.gm.ensureIndex({
  location: "2dsphere"
});