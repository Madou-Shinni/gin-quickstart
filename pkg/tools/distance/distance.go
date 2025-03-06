package distance

import "math"

const (
	earthRadius = 6371 // 地球平均半径，单位千米
)

// CalculateDistance 根据经纬度计算两地之间的距离
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := toRadians(lat1)
	φ2 := toRadians(lat2)
	Δφ := toRadians(lat2 - lat1)
	Δλ := toRadians(lon2 - lon1)

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// GetNearbyBoundingBox 根据给定经纬度和范围（公里），计算附近范围的经纬度
// 注意：不可靠计算！差距会随距离范围增大而增大
func GetNearbyBoundingBox(lat, lon, distance float64) (minLat, maxLat, minLon, maxLon float64) {
	// 计算给定经纬度范围内的矩形框
	// 计算纬度范围
	latChange := distance / earthRadius * (180 / math.Pi)
	minLat = lat - latChange
	maxLat = lat + latChange

	// 计算经度范围
	lonChange := distance / earthRadius * (180 / math.Pi) / math.Cos(lat*math.Pi/180)
	minLon = lon - lonChange
	maxLon = lon + lonChange

	return minLat, maxLat, minLon, maxLon
}
