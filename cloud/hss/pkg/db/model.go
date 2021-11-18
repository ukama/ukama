package db

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	Name  string `gorm:"not null;type:string;uniqueIndex:orgname_id_idx_case_insensetive,expression:lower(name)"`
	Imsis []Imsi
}

// Represents record in HSS db
type Imsi struct {
	gorm.Model
	OrgID          uint `gorm:"not null"`
	Org            *Org
	Imsi           string      `gorm:"index:user_imsi_idx,unique;uniqueIndex;not null;size:15;check:imsi_checker,imsi ~ $$^\\d+$$"` // IMSI Sim ID  (International mobile subscriber identity) https://www.netmanias.com/en/post/blog/5929/lte/lte-user-identifiers-imsi-and-guti
	AuthVector     *AuthVector `gorm:"embedded;embeddedPrefix:auth_"`                                                               // auth vector
	Op             []byte      `gorm:"size:16;"`                                                                                    // Pre Shared Key. This is optional and configured in operator’s DB in Authentication center and USIM. https://www.3glteinfo.com/lte-security-architecture/
	Amf            []byte      `gorm:"size:2;"`                                                                                     // Pre Shared Key. Configured in operator’s DB in Authentication center and USIM
	Key            []byte      `gorm:"size:16;"`                                                                                    // Key from the SIM
	DefaultApnName string
	UserUuid       uuid.UUID `gorm:"index:user_imsi_idx,unique;not null;type:uuid"`
	// csgid - missed
}

// Authentication Vector
// More info: https://www.3glteinfo.com/lte-security-architecture/
type AuthVector struct {
	Token        []byte `gorm:"size:16;"` // authentication token generated with AUTN = SQN * AK || AMF || MAC. It is generated only at authentication center.
	Rand         []byte `gorm:"size:16;"` // random number
	Response     []byte // Expected response generated with input (K, RAND)->f2->XRES. It is generated only at authentication center. Corresponding parameter RES is generated at USIM.
	CipherKey    []byte `gorm:"size:16;"` //  the ciphering key generated with input (K, RAND)->f3->CK. It is generated at authentication center and USIM.
	IntegrityKey []byte `gorm:"size:16;"` // integrity key generated with input (K, RAND)->f4->IK. It is generated at authentication center and USIM.
}

type User struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"uniqueIndex;not null;type:uuid"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
