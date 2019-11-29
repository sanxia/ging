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
 * database sharding
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ISharding interface {
		GetDatabaseKey() string //database key
		GetTableKey() string    //table key

		IShardingDatabase
		IShardingTable
	}

	IShardingDatabase interface {
		GetDatabaseShardingField() string     //db sharding fiele name
		SetDatabaseShardingRoute(route int32) //set db sharding route（大于0则表示直接路由到指定库，此值优先与分片字段）
	}

	IShardingTable interface {
		GetTableShardingField() string     //table sharding fiele name
		GetTableShardingCount() int32      //table sharding count
		SetTableShardingRoute(route int32) //set table sharding route（大于0则表示直接路由到指定表，此值优先与分片字段）
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get db
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetDatabaseMap(dbKey string, dbIndex int32) *gorm.DB {
	log.Printf("ging engine GetDatabaseMap dbKey: %s, shardingIndex: %d", dbKey, dbIndex)

	setting := GetApp().GetSetting()

	databaseConnection := GetDatabaseConnection(dbKey)

	isLog := !setting.Log.IsDisabled && setting.Database.IsLog
	dbMap, err := getDatabaseConnection(databaseConnection, dbIndex, isLog)
	if err != nil {
		panic(fmt.Sprintf("ging engine connection fault: %s", err.Error()))
	}

	return dbMap
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get database connection
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetDatabaseConnection(dbKey string) DatabaseConnection {
	setting := GetApp().GetSetting()

	var dbConnection DatabaseConnection

	for _, databaseConnection := range setting.Database.Connections {
		if databaseConnection.Key == dbKey {
			dbConnection = databaseConnection
			break
		}
	}

	return dbConnection
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get database connection
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func getDatabaseConnection(databaseConnection DatabaseConnection, dbIndex int32, isLog bool) (*gorm.DB, error) {
	var dbServer DatabaseServer
	for _, server := range databaseConnection.Servers {
		if server.Index == dbIndex {
			dbServer = server
		}
	}

	//dbname
	dbName := databaseConnection.Database
	if databaseConnection.ShardingCount > 0 {
		dbName = fmt.Sprintf("%s-%d", databaseConnection.Database, dbIndex)
	}

	dsn := dbServer.Username + ":" + dbServer.Password + "@tcp(" + dbServer.Host + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	dbMap, err := gorm.Open(databaseConnection.Dialect, dsn)
	if err != nil {
		log.Printf("ging engine connecting error: %s", err.Error())
	}

	dbMap.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	dbMap.DB().SetMaxIdleConns(32)
	dbMap.DB().SetMaxOpenConns(512)
	dbMap.DB().SetConnMaxLifetime(10 * time.Second)
	dbMap.LogMode(isLog)

	return dbMap, err
}
