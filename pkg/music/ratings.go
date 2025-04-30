package music

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

const (
	Zero    = Rating(0)
	OneStar = Rating(iota)
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)

type Rating int

func (r *Rating) Get() interface{} {
	return r
}

func (r *Rating) Set(value string) error {
	val, err := strconv.Atoi(value)
	rating := Rating(val)
	if err != nil || rating < Zero || rating > FiveStar {
		return fmt.Errorf("allowed values are from 0 to 5")
	} else {
		*r = rating
	}
	return nil
}

func (r Rating) String() string {
	return strconv.Itoa(int(r))
}

func (r Rating) ToCliGeneric() cli.Generic {
	return &r
}
