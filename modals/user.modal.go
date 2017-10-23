package modals
type User struct {
	Username       string `firestore:"username,omitempty"`
	Password      string `firestore:"password,omitempty"`
}
