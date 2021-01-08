package enums

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	Ratings
	Added
	Modified
	PlayCount
)
*/
type SyncType int

