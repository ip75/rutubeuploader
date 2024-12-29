package rutube

const (
	TokenUrl  = "https://rutube.ru/api/accounts/token_auth/"
	UploadUrl = "https://rutube.ru/api/video/"
)

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
