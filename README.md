# 智能骑手分配系统

## 需求规则

骑手推送规则

手动选择面试点 → 直接分配到选定面试点

未选择面试点，提供区域信息 → 根据区域内最近面试点 进行匹配

未选择面试点，也未提供区域信息 → 按面试点负载均衡（均匀分配骑手）

## 数据模型

```go
// 面试点模型
type InterviewPoint struct {
	model.Model
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`    // 纬度
	Longitude  float64 `json:"longitude"`   // 经度
	RiderCount int     `json:"rider_count"` // 当前骑手数量
}

// 骑手简历模型
type Resume struct {
    model.Model
}

// 简历推送记录模型
type ResumeLog struct {
    model.Model
    RID uint `json:"rid"` // 简历
    IID uint `json:"iid"` // 面试点
}
```

## 实现

```go
func (s *ResumeLogService) MatchIp(ctx context.Context, req ResumePushReq) (*domain.InterviewPoint, error) {
	// 规则1：直接选择站点
	if req.SelectedIID != nil {
		if ip, err := validateIp(ctx, *req.SelectedIID); err != nil {
			return nil, err
		} else {
			return ip, nil
		}
	}

	// 规则2：区域匹配
	if req.Latitude != 0 && req.Longitude != 0 {
		if req.NearLimit == 0 {
			// 默认附近20公里
			req.NearLimit = 20
		}
		if ip, err := locateRegion(ctx, req.Latitude, req.Longitude, req.NearLimit); err != nil {
			return nil, err
		} else {
			return ip, nil
		}
	}

	// 规则3：均衡分配
	return balancedAssignment(ctx)
}
```

```go
func validateIp(ctx context.Context, iid uint) (*domain.InterviewPoint, error) {
	var ip *domain.InterviewPoint
	if err := global.DB.WithContext(ctx).First(&ip, iid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("面点不存在")
		}
		return nil, err
	}
	return ip, nil
}
```

```go
func locateRegion(ctx context.Context, lat, lng float64, nearLimit float64) (*domain.InterviewPoint, error) {
	var ipList []*domain.InterviewPoint
	minLat, maxLat, minLon, maxLon := distance.GetNearbyBoundingBox(lat, lng, nearLimit)
	db := global.DB.WithContext(ctx).Model(&domain.InterviewPoint{})
	// 计算查询距离
	db = db.Select("*", fmt.Sprintf(`6371 * ACOS(
        COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) +
        SIN(RADIANS(%f)) * SIN(RADIANS(latitude))
    ) AS distance`, lat, lng, lat))
	// 粗略删选
	db = db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?", minLat, maxLat, minLon, maxLon)
	err := db.Having("distance <= ?", nearLimit). // 精确筛选（使用别名）xx km内 注意：这里Having条件放在Count查询之后 否则会报错
							Order("distance ASC").
							Order("rider_count ASC").
							Limit(1).
							Find(&ipList).Error
	if len(ipList) == 0 {
		return nil, errors.New("未找到匹配的面试点")
	}
	return ipList[0], err
}
```

```go
func balancedAssignment(ctx context.Context) (*domain.InterviewPoint, error) {
	var ips []*domain.InterviewPoint
	if err := global.DB.WithContext(ctx).Order("rider_count ASC").Limit(1).Find(&ips).Error; err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, errors.New("未找到匹配的面试点")
	}

	return ips[0], nil
}
```

## 测试用例

1. 运行测试用例TestPush

[resume_log_test.go](test%2Fresume_log_test.go)