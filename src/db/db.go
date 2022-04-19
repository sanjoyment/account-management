package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//var Database *gorm.Pool
//var err error

//type txnType string
//const (
//	DR  txnType = "dr"
//	CR txnType = "cr"
//)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "admin"
	DBNAME   = "account_mgmt"
)

type Account struct {
	Id    int    `sql:"primaryKey",json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type User struct {
	Id        int    `sql:"primaryKey",json:"id"`
	AccountId int    `json:"account_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type Txn struct {
	Id        int       `sql:"primaryKey",json:"id"`
	AccountId int       `json:"account_id"`
	UserId    int       `json:"user_id"`
	Type      string    `json:"type"`
	Detail    string    `json:"detail"`
	Amount    float32   `json:"amount"`
	Date      time.Time `json:"date"`
}

type DatabaseStore struct {
	Pool   *sql.DB
	DBName string
}

type Helpers interface {
	AutoMigrate() error
	GetDB() *DatabaseStore
}

func Connect() Helpers {
	dbinfo := fmt.Sprintf("host= %s port= %d user= %s password= %s dbname= %s sslmode=disable", HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()
	return &DatabaseStore{
		Pool:   db,
		DBName: "pg",
	}
}

func (db *DatabaseStore) GetDB() *DatabaseStore {
	return db
}

func (db *DatabaseStore) AutoMigrate() error {
	checkAccTbl := "CREATE TABLE IF NOT EXISTS accounts (" +
		"id integer NOT NULL DEFAULT nextval('accounts_id_seq'::regclass)," +
		"name text COLLATE pg_catalog.\"default\"," +
		"phone text COLLATE pg_catalog.\"default\"," +
		"CONSTRAINT accounts_pkey PRIMARY KEY (id))" +
		"TABLESPACE pg_default;" +
		"ALTER TABLE IF EXISTS public.accounts OWNER to postgres;"
	_, accErr := db.Pool.Query(checkAccTbl)

	if accErr != nil {
		return accErr
	}

	checkUsrTbl := "CREATE TABLE IF NOT EXISTS public.users(" +
		"id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass)," +
		"account_id integer," +
		"name text COLLATE pg_catalog.\"default\"," +
		"phone text COLLATE pg_catalog.\"default\"," +
		"CONSTRAINT users_pkey PRIMARY KEY (id))" +
		"TABLESPACE pg_default;" +
		"ALTER TABLE IF EXISTS public.users OWNER to postgres;"
	_, usrErr := db.Pool.Query(checkUsrTbl)

	if usrErr != nil {
		return usrErr
	}

	checkTxnTbl := "CREATE TABLE IF NOT EXISTS public.txns(" +
		"id integer NOT NULL DEFAULT nextval('txns_id_seq'::regclass)," +
		"account_id integer," +
		"user_id integer," +
		"type text COLLATE pg_catalog.\"default\"," +
		"detail text COLLATE pg_catalog.\"default\"," +
		"amount numeric," +
		"date timestamp with time zone," +
		"CONSTRAINT txns_pkey PRIMARY KEY (id))" +
		"TABLESPACE pg_default;" +
		"ALTER TABLE IF EXISTS public.txns OWNER to postgres;"
	_, txnErr := db.Pool.Query(checkTxnTbl)

	if txnErr != nil {
		return txnErr
	}
	return nil
}
