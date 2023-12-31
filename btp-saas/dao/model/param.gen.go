// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameParam = "param"

// Param mapped from table <param>
type Param struct {
	ID     int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	K      string  `gorm:"column:k;not null" json:"k"`
	V1     *string `gorm:"column:v1" json:"v1"`
	V2     *string `gorm:"column:v2" json:"v2"`
	V3     *string `gorm:"column:v3" json:"v3"`
	V4     *string `gorm:"column:v4" json:"v4"`
	V5     *string `gorm:"column:v5" json:"v5"`
	V6     *string `gorm:"column:v6" json:"v6"`
	Remark *string `gorm:"column:remark;comment:备注" json:"remark"` // 备注
}

// TableName Param's table name
func (*Param) TableName() string {
	return TableNameParam
}
