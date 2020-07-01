package request

// PostLogin - структура запроса на аунтификацию
type PostLogin struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code"` // code и state необходим для авторизации через соц. сети
	State    string `form:"state"`
}
