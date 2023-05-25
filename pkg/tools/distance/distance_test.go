package distance

import (
	"fmt"
	"testing"
)

func TestDistance(t *testing.T) {
	// SELECT *, st_distance_sphere(point(120.78650940869139,29.41920061910733),point(116.481488,39.990464)) as juli FROM `events` ORDER BY juli ASC
	distance := CalculateDistance(29.41920061910733, 120.78650940869139, 39.990464, 116.481488)
	fmt.Printf("distance :%.1f", distance)
}
