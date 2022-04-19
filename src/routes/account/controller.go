package account

import (
	"account-mgmt/src/db"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	Store db.Helpers
	from  string
}

func AccEndpoint(d db.Helpers, f string) *Endpoint {
	return &Endpoint{
		Store: d,
		from:  f,
	}
}

type ReqHome struct {
	Name  string `json:"name" binding:"required,min=1"`
	Email string `json:"email" binding:"required,min=1"`
}

type NewAccountReq struct {
	Name  string `json:"name" binding:"required,min=1"`
	Phone string `json:"phone" binding:"min=1"`
}

type UpdateAccountReq struct {
	Id    int    `json:"id" binding:"required,min=1"`
	Name  string `json:"name" binding:"min=1"`
	Phone string `json:"phone" binding:"min=1"`
}

type DeleteAccountReq struct {
	Id int `json:"id" binding:"required,min=1"`
}

func (e *Endpoint) GetAccounts(c *gin.Context) {
	var accounts []db.Account

	//log.Println("e.db--", e.Store.DBName)

	rows, err := e.Store.GetDB().Pool.Query("SELECT * FROM accounts;")

	CheckError(err)

	for rows.Next() {
		var id int
		var name string
		var phone string

		err = rows.Scan(&id, &name, &phone)

		CheckError(err)

		accounts = append(accounts, db.Account{Id: id, Name: name, Phone: phone})
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    accounts,
	})
}

func (e *Endpoint) GetAccountById(c *gin.Context) {

	accId := c.Param("id")
	var account db.Account

	err := e.Store.GetDB().Pool.QueryRow("SELECT * FROM accounts WHERE id = $1;", accId).Scan(&account.Id, &account.Name, &account.Phone)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    account,
	})
}

func (e *Endpoint) NewAccount(c *gin.Context) {
	var arg NewAccountReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}

	var lastInsertID int
	err := e.Store.GetDB().Pool.QueryRow("INSERT INTO accounts(name, phone) VALUES($1, $2) returning id;", arg.Name, arg.Phone).Scan(&lastInsertID)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Account Added Successfully, New ID: %d", lastInsertID),
	})
}

func (e *Endpoint) UpdateAccount(c *gin.Context) {
	var arg UpdateAccountReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)

	_, err1 := e.Store.GetDB().Pool.Exec("UPDATE accounts SET name = $1, phone = $2 WHERE id = $3;", arg.Name, arg.Phone, arg.Id)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Account Updated Successfully for ID: %s", arg.Id),
	})
}

func (e *Endpoint) DeleteAccount(c *gin.Context) {
	var arg DeleteAccountReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)

	_, err1 := e.Store.GetDB().Pool.Exec("DELETE FROM accounts WHERE id = $1;", arg.Id)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Account Deleted Successfully for ID: %s", arg.Id),
	})
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//--------------------

func QueryString(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")

	c.JSON(200, gin.H{
		"Name": name,
		"Age":  age,
	})
}

func PathParameters(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")

	c.JSON(200, gin.H{
		"Name": name,
		"Age":  age,
	})
}
