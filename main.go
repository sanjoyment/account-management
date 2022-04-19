package main

import (
	"account-mgmt/src/db"
	"account-mgmt/src/routes"
	"account-mgmt/src/routes/account"
	"account-mgmt/src/routes/txn"
	"account-mgmt/src/routes/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	database := db.Connect()

	if err := database.AutoMigrate(); err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println("Hello Sanjay, Welcome to Account-Management tool.")

	r := gin.Default()
	r.GET("/", routes.HomePage)
	r.POST("/", routes.HomePage)

	accEndpoint := account.AccEndpoint(database, "account")
	r.GET("/accounts", accEndpoint.GetAccounts)
	r.GET("/account/:id", accEndpoint.GetAccountById)
	r.POST("/account", accEndpoint.NewAccount)
	r.PUT("/account", accEndpoint.UpdateAccount)
	r.DELETE("/account", accEndpoint.DeleteAccount)

	usrEndpoint := user.UsrEndpoint(database, "user")
	r.GET("/users", usrEndpoint.GetUsers)
	r.GET("/user/:id", usrEndpoint.GetUserById)
	r.POST("/user", usrEndpoint.NewUser)
	r.PUT("/user", usrEndpoint.UpdateUser)
	r.DELETE("/user", usrEndpoint.DeleteUser)

	txEndpoint := txn.TxEndpoint(database, "txn")
	r.GET("/txn/:id", txEndpoint.GetTxnById)
	r.POST("/txn", txEndpoint.NewTxn)
	r.POST("/txns", txEndpoint.GetTxns)
	r.PUT("/txn", txEndpoint.UpdateTxn)
	r.DELETE("/txn", txEndpoint.DeleteTxn)

	// r.OPTIONS("/", OptionsHomePage)
	//r.GET("/query", account.QueryString)              // /query?name=sanjay&age=28
	//r.GET("/path/:name/:age", account.PathParameters) // /path/sanjay/28
	r.Run()
}
