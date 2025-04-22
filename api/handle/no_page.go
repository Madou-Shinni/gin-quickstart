package handle

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type NoPageReq struct {
	Table              string `json:"table" form:"table" uri:"table"`                                        // 表名
	Keyword            string `json:"keyword" form:"keyword" uri:"keyword"`                                  // 关键词
	Select             string `json:"select" form:"select" uri:"select"`                                     // 查询字段 逗号隔开
	FindConditionsCols string `json:"findConditionsCols" form:"findConditionsCols" uri:"findConditionsCols"` // 查询条件字段 逗号隔开
}

// NoPage 不分页查询
// @Tags     NoPage
// @Summary  查询NoPage
// @accept   application/json
// @Produce  application/json
// @Security ApiKeyAuth
// @Param    data query     NoPageReq true "不分页查询"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /noPage/{table} [get]
func NoPage(c *gin.Context) {
	var (
		req     NoPageReq
		selects = []string{"id", "name"}
		ctx     = c.Request.Context()
		count   int64
		list    = make([]map[string]interface{}, 0)
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		response.Error(c, constant.CODE_INVALID_PARAMETER, err.Error())
		return
	}

	db := global.DB.WithContext(ctx).Table(req.Table)

	req.Table = c.Param("table")
	if req.Select != "" {
		selects = strings.Split(req.Select, ",")
	}
	if req.FindConditionsCols != "" {
		slices := strings.Split(req.FindConditionsCols, ",")
		for _, v := range slices {
			db = db.Where(fmt.Sprintf("%s LIKE ?", v), "%"+req.Keyword+"%")
		}
	}

	err := db.Table(req.Table).Select(selects).Count(&count).Find(&list).Error
	if err != nil {
		response.Error(c, constant.CODE_ERR_BUSY, err.Error())
		return
	}

	response.Success(c, response.PageResponse{
		Total: count,
		List:  list,
	})
}
