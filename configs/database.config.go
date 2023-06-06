package configs

import (
	"apz-vas/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

// ConnectDb connects to the database
var DB *gorm.DB

var Services = []models.VASService{
	{
		Name:        "Mobile Airtime",
		Description: "This service allows you to buy mobile airtime",
	},
	{
		Name:        "Mobile Data",
		Description: "This service allows you to buy mobile data bundles",
	},
}

var Providers = []models.VASProvider{
	{
		Name:        "Blue Label",
		Description: "JSE-listed company that sells innovative technology for mobile commerce to emerging markets in South Africa and abroad. ",
	},
}

func ConnectDb() (*gorm.DB, error) {

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
			fmt.Println("Connecting to database...")
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
			migrate(db)
			DB = db
			return db, nil
		} else {
			return nil, err
		}
	}
	fmt.Println("Connected to database successfully")
	migrate(db)
	DB = db
	return db, nil
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.VASService{},
		&models.MobileData{},
		&models.MobileAirtime{},
		&models.ProviderService{},
		&models.SubScribedServices{},
		&models.User{},
		&models.Customer{},
		&models.VASProvider{},
	)

}
