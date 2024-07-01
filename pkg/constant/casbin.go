package constant

import "fmt"

const (
	CasbinUser = "user:%d"
	CasbinRole = "role:%d"
)

func GetCasbinUserKey(id uint) string {
	return fmt.Sprintf(CasbinUser, id)
}

func GetCasbinRoleKey(id uint) string {
	return fmt.Sprintf(CasbinRole, id)
}
