package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type FeatureCollection struct {
	Type     string
	Features []Feature
}

type Feature struct {
	Type       string
	Properties Properties
	Geometry   Geometry
}

type Properties struct {
	Id string
}

type Geometry struct {
	Type        string
	Coordinates [][]float64
}

func getCoordinates(x1, y1, x2, y2, x3, y3 float64) (bool, float64, float64) {
	xx := x2 - x1
	yy := y2 - y1
	temp := ((xx * (x3 - x1)) + (yy * (y3 - y1))) / ((xx * xx) + (yy * yy))
	X4 := x1 + xx*temp
	Y4 := y1 + yy*temp

	is_ok := true

	if X4 <= math.Max(x1, x2) && X4 >= math.Min(x1, x2) && Y4 <= math.Max(y1, y2) && Y4 >= math.Min(y1, y2) {
		is_ok = true
	} else {
		is_ok = false
	}

	return is_ok, X4, Y4
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	r := 6371.000000000
	dLat := deg2rad(lat2 - lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := r * c * 1000

	return d
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func main() {

	data, err1 := os.Open("E:/links.geojson")
	if err1 != nil {
		fmt.Println(err1)
	}
	rawData, _ := ioutil.ReadAll(data)

	var gis FeatureCollection
	err := json.Unmarshal([]byte(rawData), &gis)

	if err != nil {
		fmt.Println(2345235)
	}

	x3 := 127.027268062
	y3 := 37.499212063

	var ansLng float64
	var ansLat float64
	minDist := math.Inf(1)

	for idx := range gis.Features {
		for i := 0; i < len(gis.Features[idx].Geometry.Coordinates)-1; i++ {

			x1 := gis.Features[idx].Geometry.Coordinates[i][0]
			y1 := gis.Features[idx].Geometry.Coordinates[i][1]
			x2 := gis.Features[idx].Geometry.Coordinates[i+1][0]
			y2 := gis.Features[idx].Geometry.Coordinates[i+1][1]
			is_ok, x4, y4 := getCoordinates(x1, y1, x2, y2, x3, y3)
			if is_ok {

				dist := getDistance(y3, x3, y4, x4)
				if dist < minDist {
					minDist = dist
					ansLng = x4
					ansLat = y4
				}
			}
		}
	}

	fmt.Println(minDist, ansLng, ansLat)

}
