package entities

type User struct {
	UserID       int64 `gorm:"primaryKey;index:planned_index"`
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string

	Clan     string
	Nickname string

	HasSudo bool
}
