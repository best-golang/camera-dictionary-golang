package modals
type Set struct {
	Id       string `firestore:"-"`
	UserId       string `firestore:"user_id,omitempty"`
	Name      string `firestore:"name,omitempty"`
}