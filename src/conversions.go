package main

import (
	"strconv"
)

//ParseFloat32 since strconv.Parsefloat only does 64
//should this be embeeded?
func ParseFloat32(str string) (float32, error) {
	f64, err := strconv.ParseFloat(str, 32)
	Check(err)
	return float32(f64), err
}
