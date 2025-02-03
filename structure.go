package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

type Comment struct {
	User    int64  `json:"-" gorm:"column:id_user"`
	Post    int64  `json:"post" gorm:"column:id_post"`
	Message string `json:"message" gorm:"column:comment"`
}

func (c Comment) isValid() bool {
	if c.Post == 0 || len(c.Message) < 5 {
		return false
	}
	return true
}

func dellUser(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ROLE != 1 {
			log.Println("Попытка обращения к защищённому маршруту.")
			c.JSON(http.StatusForbidden, gin.H{"error": "restricted"})
		} else {
			id := c.Param("id")
			delUser, err := strconv.Atoi(id)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				result := db.Table(userTab).Where("id =?", delUser).Delete(&User{})
				if result.Error != nil {
					log.Println(result.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
				}
				c.JSON(http.StatusGone, gin.H{"response": "Deleted"})
			}
		}
	}
}
