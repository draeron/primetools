package enums

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	Duplicate
	Missing
)
*/
type FixType int
