package configs

import (
	"apz-vas/models"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDb connects to the database
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
	dsn := "host=" + DBHOST + " user=" + DBUSER + " password=" + DBPASS + " dbname=" + DBNAME + " port=" + DBPORT + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db.AutoMigrate(&models.Organization{})
	db.AutoMigrate(&models.VASService{})
	// db.AutoMigrate(&models.Admin{})
	// db.AutoMigrate(&models.Customer{})
	db.AutoMigrate(&models.VASProvider{})

	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}