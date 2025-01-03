package user

type User struct {
	UserName string `json:"user_name" bson:"user_name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
