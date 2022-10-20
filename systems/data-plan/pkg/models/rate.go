package models

type Rate struct {
	Id                int64            `json:"id" gorm:"primaryKey"`
	Country      string           `json:"country"`
    Country_on_cronus string           `json:"country_on_cronus"`
    Network string           `json:"network"`
    Network_id_on_cronus string           `json:"network_id_on_cronus"`
    Vpmn  string           `json:"vpmn"`
    Imsi string           `json:"imsi"`
    Sms_mo string           `json:"sms_mo"`
    Sms_mt string           `json:"sms_mt"`
    Data string           `json:"data"`
    X2g string           `json:"_2g"`
    X3g string           `json:"_3g"`
    Lte string           `json:"lte"`
    Lte_m string           `json:"lte_m"`
    Apn string           `json:"apn"`
    Created_at string           `json:"created_at"`
    Effective_at string           `json:"effective_at"`
    End_at string           `json:"end_at"`
}
