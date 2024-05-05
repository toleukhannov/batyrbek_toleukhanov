package models

import(
	"time"

)

type User struct{
	ID int
	First_Name string
	Last_Name string
	Password string
	Email string
	Phone string
	Token string
	Refresh_token string 
	Created_At time.Time
	Updated_At time.Time
}