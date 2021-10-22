package enums

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

// todo: add Eventually these
//	Traktor
//  Serato

/*
ENUM(
	ITunes
	PRIME
	File
	Rekordbox
	EngineDJ
)
*/
type LibraryType int

func (s *LibraryType) Set(value string) error {
	val, err := ParseLibraryType(value)
	if err != nil {
		return fmt.Errorf("allowed values are [%v]", strings.Join(LibraryTypeNames(), ","))
	} else {
		*s = val
	}
	return nil
}

func (s LibraryType) ToCliGeneric() cli.Generic {
	return &s
}
