package ging

import (
	"fmt"
	"log"
	"time"
)

import (
	"github.com/jinzhu/gorm"
)

/* ================================================================================
 * 数据库分片数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	//分片接口
	ISharing interface {
		GetDatabaseKey() string //数据库key
		GetTableKey() string    //表key

		ISharingDatabase
		ISharingTable
	}

	//分库接口
	ISharingDatabase interface {
		GetDatabaseShardingField() string         //获取数据库分片字段名
		GetDatabaseShardingCount() int32          //获取数据库分片数
		SetDatabaseShardingRoute(routeNode int32) //设置数据库分片路由（大于-1则表示直接路由到指定库，此值优先与分片字段）
	}

	//分表接口
	ISharingTable interface {
		GetTableShardingField() string         //获取表分片字段名
		GetTableShardingCount() int32          //获取表分片数
		SetTableShardingRoute(routeNode int32) //设置数据表分片路由（大于-1则表示直接路由到指定表，此值优先与分片字段）
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取DatabaseMap
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetDatabaseMap(dbKey string, settings Settings) *gorm.DB {
	log.Printf("Sharing GetDatabaseMap dbKey: %s", dbKey)

	var currentDatabase *DatabaseConnectionOption
	for _, database := range settings.Database.Connections {
		if database.Key == dbKey {
			currentDatabase = database
			break
		}
	}

	isLog := settings.Log.IsEnabled && settings.Database.IsLog
	dbMap, err := GetDatabaseConnection(*currentDatabase, isLog)
	if err != nil {
		panic(fmt.Sprintf("database connection fault: %s", err.Error()))
	}

	return dbMap
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 数据库链接
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetDatabaseConnection(connectionOption DatabaseConnectionOption, isLog bool) (*gorm.DB, error) {
	dsn := connectionOption.Username + ":" + connectionOption.Password + "@tcp(" + connectionOption.Host + ")/" + connectionOption.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	dbMap, err := gorm.Open(connectionOption.Dialect, dsn)

	if err != nil {
		log.Printf("Error connecting to db: %s", err.Error())
	}

	dbMap.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	dbMap.DB().SetMaxIdleConns(16)
	dbMap.DB().SetMaxOpenConns(512)
	dbMap.DB().SetConnMaxLifetime(time.Hour)
	dbMap.LogMode(isLog)

	//return &dbMap, err
	return dbMap, err
}
