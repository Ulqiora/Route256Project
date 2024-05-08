package repository

type ContactDetailDTO struct {
	Type   string `json:"type" db:"type"`
	Detail string `json:"detail" db:"detail"`
}
