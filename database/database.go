package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DbOperator struct {
	cfg   *Config
	db    *gorm.DB
	sqlDb *sql.DB
}

func (d *DbOperator) CreateDb() {
	db := "CREATE DATABASE IF NOT EXISTS " + d.cfg.DatabaseName + "DEFAULT CHARSET utf8mb4;"
	if err := d.db.Exec(db).Error; err != nil {
		panic("create database failed, err:" + err.Error())
	}
	info := fmt.Sprintf("create database: %s succeed!", d.cfg.DatabaseName)
	fmt.Println(info)
}

func (d *DbOperator) DropDb() {
	db := "DROP DATABASE IF EXISTS " + d.cfg.DatabaseName + ";"
	if err := d.db.Exec(db).Error; err != nil {
		panic("drop database failed, err:" + err.Error())
	}
	info := fmt.Sprintf("drop database: %s succeed!", d.cfg.DatabaseName)
	fmt.Println(info)
}

// ConnDb connect to db according to config,version only support mysql
func ConnDb(cfg *Config) (*DbOperator, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Url, cfg.DatabaseName)
	dialector := mysql.Open(dsn)
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanUp := func() {
		// no handler err
		_ = sqlDb.Close()
	}
	// test ping
	if err := sqlDb.Ping(); err != nil {
		return nil, cleanUp, err
	}
	// set config
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	dur, err := time.ParseDuration(cfg.ConnMaxLifeTime)
	if err != nil {
		return nil, cleanUp, err
	}
	sqlDb.SetConnMaxLifetime(dur)
	return &DbOperator{
		cfg:   cfg,
		sqlDb: sqlDb,
		db:    db,
	}, cleanUp, nil
}

// AutoHandleDB can handle db, such as create,drop,migrate
func AutoHandleDB(cfg *Config, ope string, tables ...interface{}) {
	db, cleanUp, err := ConnDb(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	switch ope {
	case "create":
		db.CreateDb()
	case "drop":
		db.DropDb()
	case "migrate":
		if err := db.db.AutoMigrate(tables...); err != nil {
			panic("auto migrate fail, err:" + err.Error())
		}
	default:
		panic("no support operation")
	}
}
