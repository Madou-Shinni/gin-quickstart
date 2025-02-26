package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"strings"
)

func importDemo(ctx context.Context, payload domain.DataImport) error {
	var (
		list          []domain.Demo
		failedReasons []domain.FailedReason
		err           error
		successCount  uint
		failureCount  uint
	)
	err = json.Unmarshal(payload.Data, &list)
	if err != nil {
		return err
	}

	return global.DB.Tx(ctx, func(ctx context.Context) error {
		for i, v := range list {
			err = func() error {
				if strings.Contains(v.Name, "error") {
					// 模拟错误
					return errors.New("名称异常")
				}
				return global.DB.WithContext(ctx).Create(&list[i]).Error
			}()

			// 统一错误处理
			if err != nil {
				failureCount++
				failedReasons = append(failedReasons, domain.FailedReason{
					Row:    fmt.Sprintf("第 %d 行", i+1),
					Reason: err.Error(),
				})
			} else {
				successCount++
			}
		}

		// 修改导入信息
		payload.SuccessCount = successCount
		payload.FailureCount = failureCount
		payload.Count = successCount + failureCount
		payload.FailedReasons = failedReasons
		if failedReasons != nil {
			payload.Status = constants.DataImportStatusFailed
		} else {
			payload.Status = constants.DataImportStatusSuccess
		}
		return global.DB.WithContext(ctx).Model(&payload).Updates(&payload).Error
	})
}
