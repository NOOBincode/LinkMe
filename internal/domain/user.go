package domain

import "time"

type User struct {
	ID         int64
	Phone      *string
	Email      string
	Nickname   string
	Password   string
	Birthday   *time.Time
	CreateTime int64
	About      string
	Deleted    bool
	Profile    Profile
}
type Profile struct {
	ID       int64
	UserID   int64
	NickName string
	Avatar   string
	Bio      string
}
