package main

type validator interface {
	isValid() bool
}

func validate(v validator) bool {
	return v.isValid()
}

type Register struct {
	Login    string `json:"login" gorm:"column:name"`
	Password string `json:"password" gorm:"password:"`
}

func (r Register) isValid() bool {
	if len(r.Login) < 5 || len(r.Password) < 5 {
		return false
	}
	return true
}

type User struct {
	Id       int64  `json:"-" gorm:"column:id"`
	Role     int8   `json:"-" gorm:"column:code_role"`
	Login    string `json:"login" gorm:"column:name"`
	Password string `json:"password" gorm:"column:password"`
}

func (u User) isValid() bool {
	if u.Id <= 0 || u.Role < 0 || len(u.Login) < 5 || len(u.Password) < 5 {
		return false
	}
	return true
}

type Post struct {
	Id      int64  `json:"-" gorm:"column:id"`
	User    int64  `json:"-" gorm:"column:id_user"`
	Content string `json:"content" gorm:"column:message"`
}

func (p Post) isValid() bool {
	if len(p.Content) == 0 {
		return false
	}
	return true
}
