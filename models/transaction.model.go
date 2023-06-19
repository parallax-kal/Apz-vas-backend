package models

import (
	"github.com/google/uuid"
)

type Topup struct {
	ID                   uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WalletId             uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:wallet_id;references:ID" json:"wallet_id"`
	Amount               float64   `gorm:"<-;not null;type:float;" json:"amount"`
	AuthCode             string    `gorm:"<-;not null;type:varchar(255);" json:"auth_code"`
	BankName             string    `gorm:"<-;not null;type:varchar(255);" json:"bank_name"`
	CardBin              string    `gorm:"<-;not null;type:varchar(255);" json:"card_bin"`
	CardLast4            string    `gorm:"<-;not null;type:varchar(255);" json:"card_last4"`
	CardName             string    `gorm:"<-;not null;type:varchar(255);" json:"card_name"`
	CardPhone            string    `gorm:"<-;not null;type:varchar(255);" json:"card_phone"`
	CardExpires          int64     `gorm:"<-;not null;type:int;" json:"card_expires"`
	CardType             string    `gorm:"<-;not null;type:varchar(255);" json:"card_type"`
	Currency             string    `gorm:"<-;not null;type:varchar(255);default:ZAR;" json:"currency"`
	ErrorDescription     string    `gorm:"<-;not null;type:varchar(255);" json:"error_description"`
	GateWay              string    `gorm:"<-;not null;type:varchar(255);" json:"gate_way"`
	GateWayTransactionId string    `gorm:"<-;not null;type:varchar(255);" json:"gate_way_transaction_id"`
	PaId                 string    `gorm:"<-;not null;type:varchar(255);" json:"paid"`
	PaymentReference     string    `gorm:"<-;not null;type:varchar(255);" json:"payment_reference"`
	TopupType            string    `gorm:"<-;not null;type:varchar(255);" json:"topup_type"`
	Type                 string    `grom:"<-;not null;type:varchar(255);" json:"type"`
	SubType              string    `grom:"<-;not null;type:varchar(255);" json:"sub_type"`
	TopUpId              float64   `grom:"<-;not null;type:int;" json:"top_up_id"`
	OrganizationWalletId float64   `gorm:"<-;not null;type:int" json:"ukheshe_wallet_id"`
	Status               string    `gorm:"<-;not null;type:varchar(255);default:PENDING" json:"status"`
	CreatedAt            int64     `gorm:"<-;not null;type:int;" json:"created_at"`
	ExpiresAt            int64     `gorm:"<-;not null;type:int;" json:"expires_at"`
}

type Withdraw struct {
	ID                   uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WalletId             uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:wallet_id;references:ID" json:"wallet_id"`
	Amount               float64   `gorm:"<-;not null;type:float;" json:"amount"`
	Currency             string    `gorm:"<-;not null;type:varchar(255);default:ZAR;" json:"currency"`
	DeliveryToPhone      string    `gorm:"<-;not null;type:varchar(255);" json:"delivery_to_phone"`
	Fee                  float64   `gorm:"<-;not null;type:float;" json:"fee"`
	GateWay              string    `gorm:"<-;not null;type:varchar(255);" json:"gate_way"`
	GateWayTransactionId string    `gorm:"<-;not null;type:varchar(255);" json:"gate_way_transaction_id"`
	Location             string    `gorm:"<-;not null;type:varchar(255);" json:"location"`
	Reference            string    `gorm:"<-;not null;type:varchar(255);" json:"reference"`
	SubType              string    `gorm:"<-;not null;type:varchar(255);" json:"sub_type"`
	Token                string    `gorm:"<-;not null;type:varchar(255);" json:"token"`
	Type                 string    `gorm:"<-;not null;type:varchar(255);" json:"type"`
	WitdrawalId          float64   `gorm:"<-;not null;type:float;" json:"witdrawal_id"`
	OrganizationWalletId float64   `gorm:"<-;not null;type:float" json:"ukheshe_wallet_id"`
	AccountName          string    `gorm:"<-;not null;type:varchar(255);" json:"account_name"`
	AccountNumber        string    `gorm:"<-;not null;type:varchar(255);" json:"account_number"`
	Bank                 string    `gorm:"<-;not null;type:varchar(255);" json:"bank"`
	BankCountry          string    `gorm:"<-;not null;type:varchar(255);" json:"bank_country"`
	BranchCode           string    `gormM:"<-;not null;type:varchar(255);" json:"branch_code"`
	ErrorDescription     string    `gorm:"<-;not null;type:varchar(255);" json:"error_description"`
	Status               string    `gorm:"<-;not null;type:varchar(255);default:PENDING" json:"status"`
	ExpiresAt            int64     `gorm:"<-;not null;type:bigint;" json:"expires_at"`
	CreatedAt            int64     `gorm:"<-;not null;type:bigint;" json:"created_at"`
	UpdatedAt            int64     `gorm:"<-;not null;type:bigint;" json:"updated_at"`
}

type Transaction struct {
	ID                   uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Description          string    `gorm:"<-;not null;type:varchar(255);" json:"description"`
	ServiceData          string    `gorm:"<-;not null;type:text" json:"service_data"`
	ServiceId            uuid.UUID `gorm:"<-;not null;type:uuid;" json:"service_id"`
	Amount               float64   `gorm:"<-;not null;type:float;" json:"amount"`
	Currency             string    `gorm:"<-;not null;type:varchar(255);" json:"currency"`
	ExternalId           uuid.UUID `gorm:"<-;not null;type:uuid" json:"external_id"`
	Rebate               float64   `gorm:"<-;not null;type:float;" json:"rebate"`
	Location             string    `gorm:"<-;not null;type:varchar(255);" json:"location"`
	ApzvasWalletId       uint32    `gorm:"<-;not null;type:int" json:"other_wallet_id"`
	OrganizationWalletId float64   `gorm:"<-;not null;type:int;" json:"organization_wallet_ID"`
	WalletId             uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:wallet_id;references:ID" json:"wallet_id"`
	CreatedAt            int64     `gorm:"autoCreateTime" json:"created_at"`
}
