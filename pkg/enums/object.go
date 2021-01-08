package enums

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	Tracks
	Playlists
	Crates
)
*/
type ObjectType int

func (e *ObjectType) Set(value string) error {
	val, err := ParseObjectType(value)
	if err != nil {
		return fmt.Errorf("allowed values are [%v]", strings.Join(ObjectTypeNames(), ","))
	} else {
		*e = val
	}
	return nil
}

func (e ObjectType) ToGeneric() cli.Generic {
	return &e
}
