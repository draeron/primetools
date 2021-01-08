package enums

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

//go:generate go-enum -f=$GOFILE --marshal --names --lower --noprefix --sql

/*
ENUM(
	Auto
	Yaml
	Json
	Text
)
*/
type FormatType int

func (f *FormatType) Set(value string) error {
	val, err := ParseFormatType(value)
	if err != nil {
		return fmt.Errorf("allowed values are [%v]", strings.Join(FormatTypeNames(), ","))
	} else {
		*f = val
	}
	return nil
}

func (f FormatType) ToCliGeneric() cli.Generic {
	return &f
}
