package enums

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	Tracks
	Playlists
	Crates
)
 */
type ObjectType int

/*
ENUM(
	Ratings
	Added
	Modified
	PlayCount
)
 */
type SyncType int
