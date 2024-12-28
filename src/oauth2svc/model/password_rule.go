package model

type PasswordPattern string

var (
	MinLength        PasswordPattern = ".{%s,}"
	LowerCase        PasswordPattern = "(?=.*[a-z])"
	UpperCase        PasswordPattern = "(?=.*[A-Z])"
	Number           PasswordPattern = "(?=.*[0-9])"
	SpecialCharacter PasswordPattern = "(?=.*[@$!%*?&])"
)

type PasswordRule struct {
	ID          string          `gorm:"column:id;type:uuid;default:uuid_generated_v4();primaryKey" json:"-"`
	Name        string          `gorm:"column:name" json:"name"`
	DisplayName string          `gorm:"column:display_name" json:"display_name"`
	Pattern     PasswordPattern `gorm:"column:pattern" json:"pattern"`
}

type PasswordRules []PasswordRule
