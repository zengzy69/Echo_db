package db

import (
	"echoDB/config"
	"echoDB/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// LoadDataFromMySQL 从MySQL加载数据

func LoadDataFromMySQL(config *config.Config) ([]models.Item, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.User, config.Database.Password,
		config.Database.Host, config.Database.Port,
		config.Database.DBName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	var items []models.Item

	db.Find(&items)

	return items, nil
}
