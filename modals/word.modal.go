package modals
type Word struct {
	Id       string `firestore:"-"`
	SetId       string `firestore:"set_id,omitempty"`
	Type      string `firestore:"type,omitempty"`
	Content      string `firestore:"content,omitempty"`
	Language     string `firestore:"language,omitempty"`
	Translation      string `firestore:"translation,omitempty"`
	Pronunciation      string `firestore:"pronunciation,omitempty"`
}