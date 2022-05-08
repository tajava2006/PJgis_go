package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

// // geojson을 파싱하기 위한 struct
// 최상위 FeatureCollection. Feature의 리스트를 features 라는 키로 가지고 있음
type FeatureCollection struct {
	Type     string
	Features []Feature
}

// 각각의 geojson 객체들
type Feature struct {
	Type       string
	Properties Properties
	Geometry   Geometry
}

// geojson 객체들의 속성
type Properties struct {
	Id string
}

// geojson 객체들의 위치정보를 가지고 있음. coordinates 라는 키 안에 각 좌표의 리스트가 포함되어 있음
type Geometry struct {
	Type        string
	Coordinates [][]float64
}

// 경위도로 표현된 두 선분의 끝점(lng1,lat1,lng2,lat2)과 target(lng3,lat3)을 받는다
// target에서 선분에 내린 수선의 발의 경위도를 리턴한다
// 선분과 타겟의 사잇각이 90도가 넘으면 수선의 발이 선분에 위치 하지 못하므로
// 해당 정보를 bool 타입으로 넘겨준다
func getOrthogonalCoordinates(lng1, lat1, lng2, lat2, lng3, lat3 float64) (bool, float64, float64) {
	xx := lng2 - lng1
	yy := lat2 - lat1
	temp := ((xx * (lng3 - lng1)) + (yy * (lat3 - lat1))) / ((xx * xx) + (yy * yy))
	X4 := lng1 + xx*temp
	Y4 := lat1 + yy*temp

	is_ok := true

	if X4 <= math.Max(lng1, lng2) && X4 >= math.Min(lng1, lng2) && Y4 <= math.Max(lat1, lat2) && Y4 >= math.Min(lat1, lat2) {
		is_ok = true
	} else {
		is_ok = false
	}

	return is_ok, X4, Y4
}

// 경위도로 표현된 두 점 사이의 거리를 구한다
func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	r := 6371.000000000 // 지구의 반지름(km)
	dLat := deg2rad(lat2 - lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := r * c * 1000 // 미터로 환산하여 리턴한다

	return d
}

// 계산을 위해 각도의 단위를 라디안으로 바꾼다
func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func main() {

	// links와 target라는 플래그를 커맨드라인으로 받는다
	links := flag.String("links", "links.geojson", "geojson 데이터파일명")
	target := flag.String("target", "127.027268062,37.499212063", "target의 경위도")

	flag.Parse()

	// 쉼표로 표현된 target를 파싱하여 floa64로 변환하여 target변수에 담는다
	temp1 := *target
	temp2 := strings.Split(temp1, ",")
	lng3, _ := strconv.ParseFloat(temp2[0], 64)
	lat3, _ := strconv.ParseFloat(temp2[1], 64)

	// 주어진 geojson 파일을 연다
	data, err1 := os.Open("./" + *links)
	if err1 != nil {
		fmt.Println(err1)
	}
	rawData, _ := ioutil.ReadAll(data)

	// 위에서 선언해준 struct를 사용하여 geojson 객체(FeatureCollection)로 담는다
	var gis FeatureCollection
	err := json.Unmarshal([]byte(rawData), &gis)
	if err != nil {
		fmt.Println(2345235)
	}

	// 정답을 담을 변수들을 선언한다.
	var ansLng float64
	var ansLat float64
	minDist := math.Inf(1)

	// gis 객체 안에 있는 모든 feature를 돌면서 target 와의 거리를 측정한다
	for idx := range gis.Features {
		// 각각의 feature가 단순한 line이 아니라 lineString이기 때문에 해당 feature의 crrodinates의 인접한 두 좌표를 두 끝점으로 하는 line(=link)와 target과의 거리를 측정한다
		for i := 0; i < len(gis.Features[idx].Geometry.Coordinates)-1; i++ {

			lng1 := gis.Features[idx].Geometry.Coordinates[i][0]
			lat1 := gis.Features[idx].Geometry.Coordinates[i][1]
			lng2 := gis.Features[idx].Geometry.Coordinates[i+1][0]
			lat2 := gis.Features[idx].Geometry.Coordinates[i+1][1]
			is_ok, lng4, lat4 := getOrthogonalCoordinates(lng1, lat1, lng2, lat2, lng3, lat3)
			if is_ok {
				dist := getDistance(lat3, lng3, lat4, lng4)
				// 측정된 거리가 이제까지의 측정된 최소거리 보다 짧다면 정답을 갱신한다
				if dist < minDist {
					minDist = dist
					ansLng = lng4
					ansLat = lat4
				}
			}
		}
	}

	// 정답을 출력한다
	fmt.Println(minDist, ansLng, ansLat)

}
