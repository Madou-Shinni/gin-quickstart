package pagelimit

import (
	"fmt"
	"testing"
)

func TestPageLimit(t *testing.T) {
	offset, limit := OffsetLimit(0, 1000)

	fmt.Printf("offset:%v, limit:%v", offset, limit)
}
