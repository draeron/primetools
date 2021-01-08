package enums

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	ITunes
	PRIME
	File
	Traktor
)
*/
type SourceType int

func (s *SourceType) Set(value string) error {
	val, err := ParseSourceType(value)
	if err != nil {
		return fmt.Errorf("allowed values are [%v]", strings.Join(SourceTypeNames(), ","))
	} else {
		*s = val
	}
	return nil
}

func (s SourceType) ToCliGeneric() cli.Generic {
	return &s
}
