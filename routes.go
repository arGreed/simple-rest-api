package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	pingRoute     string = "/ping"
	registerRoute string = "/register"
	loginRoute    string = "/login"

	postRoute    string = "/new/post"
	commentRoute string = "/new/comment"

	delUserRoute    string = "/user/delete/:id"
	delPostRoute    string = "/post/delete/:id"
	delCommentRoute string = "/comment/delete/:id"

	usersRoute    string = "/user/all"
	postsRoute    string = "/user/posts"
	commentsRoute string = "/user/comments"

	userTab    string = "simple_rest_app.user"
	postTab    string = "simple_rest_app.post"
	commentTab string = "simple_rest_app.post_comment"

	//Временное решение
	USER int64
	ROLE int8
)

func logMiddleware(c *gin.Context) {
	log.Println(c.Request.Method, c.Request.URL.Path)
	c.Next()
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"connected": "ok"})
}

func register(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var register Register
		var user User
		err := c.ShouldBindJSON(&register)
		if err != nil || !validate(register) {
			log.Println("Пользователь при регистрации ввёл некорректные параметры!")
		}
		result := db.Table(userTab).Where("name = ?", register.Login).First(&user)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			log.Println("Ошибка работы с БД")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка работы с БД"})
		}
		if user.Id == 0 {
			db.Table(userTab).Create(&register)
			c.JSON(http.StatusCreated, gin.H{"success": "Пользователь успешно зарегистрирован!"})
		} else {
			log.Println("Попытка создания пользователя с совпадающим login'ом")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Попытка создания пользователя с совпадающим login'ом"})
		}
	}
}

func login(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var login Register
		var user User
		err := c.ShouldBindJSON(&login)
		if err != nil || !validate(login) {
			log.Println("Переданы данные не верного формата.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Переданы данные не верного формата."})
		}

		result := db.Table(userTab).Where("name = ? and password = ?", login.Login, login.Password).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Пользователь не найден")
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден."})
		} else if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		}

		//c.Set("id", user.Id)
		//c.Set("role", user.Role)
		USER = user.Id
		ROLE = user.Role
		c.JSON(http.StatusOK, gin.H{"Успех": "Пользователь авторизовался!"})
	}
}

func post(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var post Post
		err := c.ShouldBindJSON(&post)
		if err != nil || !validate(post) {
			log.Println("Получены некорректные данные")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Получены некорректные данные."})
		}
		//fmt.Println(c.Get("id"))
		//userIdStr, _ := c.Get("id")
		//userId := userIdStr.(int64)

		//fmt.Println(userId)
		//post.User = int64(userId)
		post.User = USER
		result := db.Table(postTab).Create(&post)

		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		} else {
			c.JSON(http.StatusCreated, gin.H{"created": "Пост успешно создан!"})
		}
	}
}

func comment(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var comment Comment
		err := c.ShouldBindJSON(&comment)
		if err != nil || !validate(comment) {
			log.Println("Переданы некорректные данные!")
			c.JSON(http.StatusBadRequest, gin.H{"errror": "Переданы некорректные данные!"})
		}
		var post Post

		result := db.Table(postTab).Where("id = ?", comment.Post).First(&post)
		if result.Error == gorm.ErrRecordNotFound {
			log.Println(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		} else if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		}
		comment.User = USER
		result = db.Table(commentTab).Create(&comment)
		if result.Error != nil {
			log.Println(result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		} else {
			c.JSON(http.StatusCreated, gin.H{"response": "Record created"})
		}
	}
}

func delUser(db *gorm.DB) func(c *gin.Context) {
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

func delPost(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		delPost, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			result := db.Table(postTab).Where("id_user =? and id = ?", USER, delPost)
			resultAdm := db.Table(postTab).Where("id =?", delPost)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			} else if ROLE == 1 && resultAdm.Error != nil {
				log.Println(resultAdm.Error)
				c.JSON(http.StatusBadRequest, gin.H{"error": resultAdm.Error})
			} else {
				result = db.Table(postTab).Where("id =?", delPost).Delete(&Post{})
				if result.Error != nil {
					log.Println(result.Error)
					c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
				}
				c.JSON(http.StatusGone, gin.H{"response": "Пост успешно удалён."})
			}
		}
	}
}

func delComment(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var comment Comment
		id := c.Param("id")
		delComment, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			resultComment := db.Table(commentTab).Where("id = ?", delComment).First(&comment)
			if resultComment.Error == gorm.ErrRecordNotFound {
				log.Println("Получены некорректные данные")
				c.JSON(http.StatusBadRequest, gin.H{"error": gorm.ErrRecordNotFound})
			} else if resultComment.Error != nil {
				log.Println(resultComment.Error)

			}

			if ROLE == 1 {
				delAdm := db.Table(commentTab).Where("id = ?", delComment).Delete(&Comment{})
				if delAdm.Error != nil {
					log.Println(delAdm.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"error": delAdm.Error})
				}
			} else {
				del := db.Table(commentTab).Where("id = ? and id_user = ?", delComment, USER).Delete(&Comment{})
				if del.Error != nil {
					log.Println(del.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"error": del.Error})
				}
			}
			c.JSON(http.StatusGone, gin.H{"response": "deleted"})
		}
	}
}

func allUsers(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var users []User
		if ROLE == 1 {
			result := db.Table(userTab).Find(&users)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			}
			c.JSON(http.StatusOK, gin.H{"users": users})
		} else {
			log.Println("Попытка обращения к защищённому маршруту!")
			c.JSON(http.StatusForbidden, gin.H{"error": "Попытка обращения к защищённому маршруту!"})
		}
	}
}

func userPosts(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var posts []Post
		if ROLE == 1 {
			result := db.Table(postTab).Find(&posts)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})
		} else {
			result := db.Table(postTab).Where("id_user = ?", USER).Find(&posts)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})
		}
	}
}

func userComments(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var comments []Comment
		if ROLE == 1 {
			result := db.Table(commentTab).Find(&comments)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			}
			c.JSON(http.StatusOK, gin.H{"comments": comments})
		} else {
			result := db.Table(postTab).Where("id_user = ?", USER).Find(&comments)
			if result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			}
			c.JSON(http.StatusOK, gin.H{"comments": comments})
		}
	}
}
