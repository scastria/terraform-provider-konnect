package client

const (
	UserPath    = "users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	Id            string `json:"id,omitempty"`
	Email         string `json:"email,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	PreferredName string `json:"preferred_name,omitempty"`
	Active        bool   `json:"active"`
}
type UserCollection struct {
	Users []User `json:"data"`
}
