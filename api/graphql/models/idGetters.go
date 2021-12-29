package models

type IDGetter interface {
	GetID() string
}

func (u *User) GetID() string {
	return u.ID
}

var _ IDGetter = &User{}

func (uu *UpdateUserInput) GetID() string {
	return uu.ID
}

var _ IDGetter = &UpdateUserInput{}
