package defs

type MongoCountResult struct {
	N  int  `bson:"n"`
	Ok bool `bson:"ok"`
}
