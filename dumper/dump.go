package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/RATDistributedSystems/utilities"
	"github.com/beevik/etree"
	"github.com/gocql/gocql"
)

var sessionGlobal *gocql.Session
var configs = utilities.GetConfigurationFile("config.json")
var page_size, _ = strconv.Atoi(configs.GetValue("cassandra_page_size"))
var time_out,_  = strconv.Atoi(configs.GetValue("cassandra_timeout"))

func main() {
	//establish global connection param


	cluster := gocql.NewCluster(configs.GetValue("cassandra_ip"))
	cluster.ConnectTimeout = time.Second * time.Duration(time_out)
	cluster.Keyspace = configs.GetValue("cassandra_keyspace")
	proto, err := strconv.Atoi(configs.GetValue("cassandra_proto"))
	if err != nil {
		panic("Cassandra protocol version not int")
	}
	cluster.ProtoVersion = proto

	session, err := cluster.CreateSession()
	sessionGlobal = session
	if err != nil {
		panic(err)
	}

	dump("1000users")

}

	
func dump(filename string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := etree.NewElement("log")
	doc.SetRoot(root)
	filename = strings.TrimSpace(strings.TrimSuffix(filename, "\n"))
	filename = addTimestampToFilename(filename)
	log.Printf("Dumping all users to %s\n", filename)

	var transactionTime string
	var server string
	var transactionNum string
	var command string
	var stockSymbol string
	var funds string
	var quoteservertime string
	var cryptokey string
	var price string
	var userId string
	//var action string
	var count int
	//var errorMessage string
	//var debugMessage string

	//check if user commands
	/*
	if err := sessionGlobal.Query("SELECT count(*) FROM usercommands").Scan(&count); err != nil {
		panic(err)
	}
	*/

	count = 1
	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands").PageSize(page_size).Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			user_command(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if quote server events
	/*
	if err := sessionGlobal.Query("SELECT count(*) FROM quote_server").Scan(&count); err != nil {
		panic(err)
	}
	*/
	count = 1
	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, quoteservertime , userid, stocksymbol, price, cryptokey FROM quote_server").PageSize(page_size).Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &quoteservertime, &userId, &stockSymbol, &price, &cryptokey) {
			quote_server(doc, transactionTime, server, transactionNum, quoteservertime, userId, stockSymbol, price, cryptokey)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	/*
	//check if account transaction events
	if err := sessionGlobal.Query("SELECT count(*) FROM account_transaction").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, action, userid, funds FROM account_transaction ").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &action, &userId, &funds) {
			account_transaction(doc, transactionTime, server, transactionNum, action, userId, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if system event
	if err := sessionGlobal.Query("SELECT count(*) FROM system_event").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM system_event ").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			system_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if error event
	if err := sessionGlobal.Query("SELECT count(*) FROM error_event").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, errorMessage FROM error_event ").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &errorMessage) {
			error_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds, errorMessage)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if debug event
	if err := sessionGlobal.Query("SELECT count(*) FROM debug_event").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, debugMessage FROM debug_event ").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &debugMessage) {
			debug_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds, debugMessage)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	*/

	doc.Indent(2)
	doc.WriteToFile(filename)
	fmt.Println("Done!")
}

func addTimestampToFilename(f string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", f, t.Format("2006-01-02_15-04-05"))
}
