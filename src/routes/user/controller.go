package user

import (
	"account-mgmt/src/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Endpoint struct {
	Store db.Helpers
	from  string
}

func UsrEndpoint(d db.Helpers, f string) *Endpoint {
	return &Endpoint{
		Store: d,
		from:  f,
	}
}

type NewUserReq struct {
	AccountId int    `json:"account_id" binding:"required,min=1"`
	Name      string `json:"name" binding:"required,min=1"`
	Phone     string `json:"phone" binding:"min=1"`
}

type UpdateUserReq struct {
	Id        int    `json:"id" binding:"required,min=1"`
	AccountId int    `json:"account_id" binding:"required,min=1"`
	Name      string `json:"name" binding:"min=1"`
	Phone     string `json:"phone" binding:"min=1"`
}

type DeleteUserReq struct {
	Id        int `json:"id" binding:"required,min=1"`
	AccountId int `json:"account_id" binding:"required,min=1"`
}

func (e *Endpoint) GetUsers(c *gin.Context) {
	var users []db.User
	rows, err := e.Store.GetDB().Pool.Query("SELECT * FROM users;")

	CheckError(err)

	for rows.Next() {
		var id int
		var accountId int
		var name string
		var phone string

		err = rows.Scan(&id, &accountId, &name, &phone)

		CheckError(err)

		users = append(users, db.User{Id: id, AccountId: accountId, Name: name, Phone: phone})
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    users,
	})
}

func (e *Endpoint) GetUserById(c *gin.Context) {

	usrId := c.Param("id")
	accId := c.Param("accountId")
	var user db.User
	err := e.Store.GetDB().Pool.QueryRow("SELECT * FROM users WHERE id = $1 AND account_id = $2;", usrId, accId).Scan(&user.Id, &user.AccountId, &user.Name, &user.Phone)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
	})
}

func (e *Endpoint) NewUser(c *gin.Context) {
	var arg NewUserReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}

	var lastInsertID int
	err := e.Store.GetDB().Pool.QueryRow("INSERT INTO users(account_id, name, phone) VALUES($1, $2, $3) returning id;", arg.AccountId, arg.Name, arg.Phone).Scan(&lastInsertID)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("User Added Successfully, New ID: %d", lastInsertID),
	})
}

func (e *Endpoint) UpdateUser(c *gin.Context) {
	var arg UpdateUserReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)
	_, err1 := e.Store.GetDB().Pool.Exec("UPDATE users SET name = $1, phone = $2 WHERE id = $3 AND account_id = $4;", arg.Name, arg.Phone, arg.Id, arg.AccountId)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("User Updated Successfully for ID: %s", arg.Id),
	})
}

func (e *Endpoint) DeleteUser(c *gin.Context) {
	var arg DeleteUserReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)
	_, err1 := e.Store.GetDB().Pool.Exec("DELETE FROM users WHERE id = $1 AND account_id = $2;", arg.Id, arg.AccountId)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("User Deleted Successfully for ID: %s", arg.Id),
	})
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
