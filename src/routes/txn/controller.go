package txn

import (
	"account-mgmt/src/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type Endpoint struct {
	Store db.Helpers
	from  string
}

func TxEndpoint(d db.Helpers, f string) *Endpoint {
	return &Endpoint{
		Store: d,
		from:  f,
	}
}

type GetTxnsReq struct {
	AccountId int `json:"account_id" binding:"required,min=1"`
	UserId    int `json:"user_id"`
}
type NewTxnReq struct {
	AccountId int       `json:"account_id" binding:"required,min=1"`
	UserId    int       `json:"user_id" binding:"required,min=1"`
	Type      string    `json:"type" binding:"required,min=1"`
	Detail    string    `json:"detail" binding:"required,min=1"`
	Amount    float32   `json:"amount" binding:"required,min=1"`
	Date      time.Time `json:"date"`
}

type UpdateTxnReq struct {
	Id        int       `json:"id" binding:"required,min=1"`
	AccountId int       `json:"account_id" binding:"required,min=1"`
	UserId    int       `json:"user_id" binding:"required,min=1"`
	Type      string    `json:"type" binding:"required,min=1"`
	Detail    string    `json:"detail" binding:"required,min=1"`
	Amount    float32   `json:"amount" binding:"required,min=1"`
	Date      time.Time `json:"date"`
}

type DeleteTxnReq struct {
	Id        int `json:"id" binding:"required,min=1"`
	AccountId int `json:"account_id" binding:"required,min=1"`
	UserId    int `json:"user_id" binding:"required,min=1"`
}

func (e *Endpoint) GetTxns(c *gin.Context) {
	var arg GetTxnsReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}

	andWhere := ""
	if arg.UserId > 0 {
		andWhere = fmt.Sprintf(" AND user_id = %d", arg.UserId)
	}

	var txns []db.Txn
	rows, err := e.Store.GetDB().Pool.Query("SELECT * FROM txns WHERE account_id = $1"+andWhere+";", arg.AccountId)

	CheckError(err)

	for rows.Next() {
		var id int
		var accountId int
		var userId int
		var type_ string
		var detail string
		var amount float32
		var date time.Time

		err = rows.Scan(&id, &accountId, &userId, &type_, &detail, &amount, &date)

		CheckError(err)

		txns = append(txns, db.Txn{Id: id, AccountId: accountId, UserId: userId, Type: type_, Detail: detail, Amount: amount, Date: date})
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    txns,
	})
}

func (e *Endpoint) GetTxnById(c *gin.Context) {
	txnId := c.Param("id")
	usrId := c.Param("userid")
	accId := c.Param("accountId")
	var txn db.Txn
	err := e.Store.GetDB().Pool.QueryRow(
		"SELECT * FROM txns WHERE id = $1 AND account_id = $2 AND user_id = $3;",
		txnId,
		accId,
		usrId).Scan(
		&txn.Id,
		&txn.AccountId,
		&txn.UserId,
		&txn.Type,
		&txn.Detail,
		&txn.Amount,
		&txn.Date)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    txn,
	})
}

func (e *Endpoint) NewTxn(c *gin.Context) {
	var arg NewTxnReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}

	var lastInsertID int
	err := e.Store.GetDB().Pool.QueryRow(
		"INSERT INTO txns(account_id, user_id, type, detail, amount, date) VALUES($1, $2, $3, $4, $5, $6) returning id;",
		arg.AccountId,
		arg.UserId,
		arg.Type,
		arg.Detail,
		arg.Amount,
		arg.Date,
	).Scan(&lastInsertID)

	CheckError(err)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Txn Added Successfully, New ID: %d", lastInsertID),
	})
}

func (e *Endpoint) UpdateTxn(c *gin.Context) {
	var arg UpdateTxnReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)
	_, err1 := e.Store.GetDB().Pool.Exec("UPDATE txns SET type = $1, detail = $2, amount = $3, date = $4 WHERE id = $5 AND account_id = $6 AND user_id = $7;",
		arg.Type,
		arg.Detail,
		arg.Amount,
		arg.Date,
		arg.Id,
		arg.AccountId,
		arg.UserId)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Txn Updated Successfully for ID: %d", arg.Id),
	})
}

func (e *Endpoint) DeleteTxn(c *gin.Context) {
	var arg DeleteTxnReq
	if err := c.BindJSON(&arg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(arg)
	_, err1 := e.Store.GetDB().Pool.Exec("DELETE FROM Txn WHERE id = $1 AND account_id = $2; AND user_id = $3;", arg.Id, arg.AccountId, arg.UserId)

	CheckError(err1)

	c.JSON(200, gin.H{
		"success": true,
		"data":    fmt.Sprintf("Txn Deleted Successfully for ID: %s", arg.Id),
	})
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
