package configs

import (
	"apz-vas/models"
	"fmt"
	"os"
	"strings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDb connects to the database
var DB *gorm.DB

var Services = []models.VASService{
	{
		Name:        "Mobile Airtime",
		NickName:    "airtime",
		Rebate:      0.05,
		Description: "This service allows you to buy mobile airtime",
	},
	{
		Name:        "Mobile Data",
		NickName:    "bundle",
		Rebate:      0.05,
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
		&models.SubscribedServices{},
		&models.User{},
		&models.Customer{},
		&models.VASProvider{},
	)

	for _, provider := range Providers {
		newProvider := db.Create(&provider)
		if newProvider.Error != nil {
			fmt.Println(newProvider.Error)
		}
	}
	
	var providers []models.VASProvider
	db.Find(&providers)
	for _, service := range Services {
		// all the services be of BlueLabel
		service.ProviderId = providers[0].ID	
		// create or update
		newService := db.Create(&service)
		if newService.Error != nil {
			fmt.Println(newService.Error)
		}

	}
	// FOR EXAMPLE, BLUE LABEL HAS MOBILE AIRTIME AND MOBILE DATA



}
