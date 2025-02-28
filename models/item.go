package models

// Item 是数据库中存储的单个数据项
type Item struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"` // 表示主键
	Value      string `gorm:"type:text"`                // 存储 Value，使用 text 类型，适应较大的值
	ExpiryTime int64  `gorm:"type:bigint"`              // 存储 ExpiryTime，使用 BIGINT 类型
}
