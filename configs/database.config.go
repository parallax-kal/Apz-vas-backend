package configs

import (
	"apz-vas/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

// ConnectDb connects to the database
var DB *gorm.DB

func ConnectDb() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	DBHOST := os.Getenv("DBHOST")
	DBPORT := os.Getenv("DBPORT")
	DBUSER := os.Getenv("DBUSER")
	DBPASS := os.Getenv("DBPASS")
	DBNAME := os.Getenv("DBNAME")
	fmt.Println("Connecting to database...")
	dsn := "host=" + DBHOST + " user=" + DBUSER + " password=" + DBPASS + " dbname=" + DBNAME + " port=" + DBPORT + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if error is about database does not exist, create the database
	if err != nil {
		if strings.Contains(err.Error(), "database \""+DBNAME+"\" does not exist") {
			fmt.Println("Database does not exist, creating database...")
			dsn = "host=" + DBHOST + " user=" + DBUSER + " password=" + DBPASS + " port=" + DBPORT + " sslmode=disable TimeZone=Asia/Shanghai"
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, err
			}
			// Create database
			err = db.Exec("CREATE DATABASE \"" + DBNAME + "\"").Error
			if err != nil {
				return nil, err
			}
			fmt.Println("Database created successfully")
			// Connect to database
			dsn = "host=" + DBHOST + " user=" + DBUSER + " password=" + DBPASS + " dbname=" + DBNAME + " port=" + DBPORT + " sslmode=disable TimeZone=Asia/Shanghai"
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, err
			}
			// enable uuid
			fmt.Println("Connected to database successfully")
			fmt.Println("Enabling uuid-ossp extension...")
			err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
			if err != nil {
				return nil, err
			}
			fmt.Println("Enabled uuid-ossp extension successfully")
			db.AutoMigrate(&models.VASService{})
			db.AutoMigrate(&models.Admin{})
			db.AutoMigrate(&models.Customer{})
			db.AutoMigrate(&models.Organization{})
			db.AutoMigrate(&models.VASProvider{})
		} else {
			return nil, err
		}
		return nil, err
	}
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database successfully")
	db.AutoMigrate(&models.VASService{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Customer{})
	db.AutoMigrate(&models.Organization{})
	db.AutoMigrate(&models.VASProvider{})
	DB = db
	return db, nil
}
