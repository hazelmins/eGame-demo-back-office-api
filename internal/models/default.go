/*
 * @Description:table 表格在這邊加入鄉親
 */
package models

func GetModels() []interface{} {
	return []interface{}{
		&AdminUsers{}, &User{}, &SuperAdmin{},
	}
}
