package models

import (
	"time"
)

type Certificate struct {
	ID                 uint      `json:"id"`
	Number             string    `json:"number"`
	OverallLength      string    `json:"overall_length"`
	BladeLength        string    `json:"blade_length"`
	HandleLength       string    `json:"handle_length"`
	BladeMaterial      string    `json:"blade_material"`
	BladeThickness     string    `json:"blade_thickness"`
	KissanLength       string    `json:"kissan_length"`
	Hamon              string    `json:"hamon"`
	Mekugi             int8      `json:"mekugi"`
	Sakihaba           string    `json:"sakihaba"`
	Motohaba           string    `json:"motohaba"`
	SayaMaterial       string    `json:"saya_material"`
	TsubaMaterial      string    `json:"tsuba_material"`
	HabakiMaterial     string    `json:"habaki_material"`
	ItoSageoMaterial   string    `json:"ito_sageo_material"`
	ForgeID            uint      `json:"forge_id"`
	SwordsmithID       uint      `json:"swordsmith_id"`
	ForgeName          string    `json:"forge_name" gorm:"->;column:forge_name"`
	SwordsmithName     string    `json:"swordsmith_name" gorm:"->;column:swordsmith_name"`
	ForgeImageURL      string    `json:"forge_image_url" gorm:"->;column:forge_image_url"`
	SwordsmithImageURL string    `json:"swordsmith_image_url" gorm:"->;column:swordsmith_image_url"`
	DateCompleted      string    `json:"date_completed"`
	No                 string    `json:"no"`
	CreatedAt          time.Time `json:"created_at"`
}

func (Certificate) TableName() string {
	return "certificates"
}
