package models

import (
	"github.com/google/uuid"
)

type Organization struct {
	ID                        uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserId                    uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:user_id;references:ID" json:"user_id"`
	APIKey                    uuid.UUID `gorm:"<-:create;unique;type:uuid;default:uuid_generate_v4()" json:"api_key"`
	Company_Name              string    `gorm:"<-;not null;type:varchar(255)" json:"company_name"`
	Owner_Name                string    `gorm:"<-;not null;type:varchar(255)" json:"owner_name"`
	Email                     string    `gorm:"<-;unique;not null;type:varchar(255)" json:"email"`
	Company_Number            string    `gorm:"<-;type:varchar(255)" json:"company_number"`
	Phone_Number1             string    `gorm:"<-;type:varchar(255)" json:"phone_number1"`
	Phone_Number2             string    `gorm:"<-;type:varchar(255)" json:"phone_number2"`
	Trading_Name              string    `gorm:"<-;type:varchar(255)" json:"trading_name"`
	Industrial_Sector         string    `gorm:"<-;type:varchar(255)" json:"industrial_sector"`
	Industrial_Classification string    `gorm:"<-;type:varchar(255)" json:"industrial_classification"`
	Tax_Number                string    `gorm:"<-;type:varchar(255)" json:"tax_number"`
	Ukheshe_Id                uint32    `gorm:"<-;type:int" json:"ukheshe_id"`
	Registration_Date         string    `gorm:"<-;type:varchar(255)" json:"registration_date"`
	Organization_Type         string    `gorm:"<-;type:varchar(255)" json:"organization_type"`
	BusinessType              string    `gorm:"<-;type:varchar(255)" json:"businesstype"`
	Status                    string    `gorm:"<-;type:varchar(255);default:Active" json:"status"`
	Version                   int       `gorm:"<-;type:float" json:"version"`
	UpdatedAt                 int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt                 int64     `gorm:"autoCreateTime" json:"created_at"`
}
