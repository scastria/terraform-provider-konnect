package client

const (
	UserPath    = "v3/users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	Id            string `json:"id,omitempty"`
	Email         string `json:"email,omitempty"`
	FullName      string `json:"full_name"`
	PreferredName string `json:"preferred_name"`
	Active        bool   `json:"active,omitempty"`
}
type UserCollection struct {
	Users []User `json:"data"`
}
