package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/beevik/etree"
	"github.com/gocql/gocql"
)

var sessionGlobal *gocql.Session
var configs = utilities.Load()

func main() {
	configs.Pause()
	hostNoSpace := strings.TrimSpace(configs.GetValue("auditdb_ip"))
	keyspace := configs.GetValue("auditdb_keyspace")
	proto := configs.GetValue("auditdb_proto")
	ratdatabase.InitCassandraConnection(hostNoSpace, keyspace, proto)
	sessionGlobal = ratdatabase.CassandraConnection

	addr, protocol := configs.GetListnerDetails("audit")
	l, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer l.Close()

	var wg sync.WaitGroup
	log.Printf("Listeniing on %s", addr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn, &wg)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	for {
		tp := textproto.NewReader(bufio.NewReader(conn))
		msg, err := tp.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Println("Encountered Error while processing connection")
				log.Println(err)
			}
			continue
		}
		go handleRequest(msg, wg)
	}
}

func handleRequest(msg string, wg *sync.WaitGroup) {
	log.Printf("Received request: %s\n", msg)
	result := strings.Split(string(msg), ",")

	//add user event to database
	//time, server, transactionNum, command, userid, funds

	if result[0] == "User" {
		wg.Add(1)
		logUserEvent(result)
		wg.Done()

	}

	if result[0] == "Quote" {
		wg.Add(1)
		logQuoteEvent(result)
		wg.Done()

	}

	if result[0] == "System" {
		wg.Add(1)
		logSystemEvent(result)
		wg.Done()

	}

	if result[0] == "Error" {
		wg.Add(1)
		logErrorEvent(result)
		wg.Done()
	}

	if result[0] == "Account" {
		wg.Add(1)
		logAccountTransactionEvent(result)
		wg.Done()

	}

	if result[0] == "Debug" {
		wg.Add(1)
		logDebugEvent(result)
		wg.Done()

	}

	if result[0] == "DUMPLOG" {
		wg.Wait()
		time.Sleep(time.Second * 10)
		if len(result) == 3 {
			//DUMP with specific user
			dumpUser(result[1], result[2])
		} else if len(result) == 2 {
			//DUMP everything
			dump(result[1])
		}
	}
}

func logUserEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stockSymbol, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "' , '" + result[7] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logQuoteEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO quote_server (time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "' , " + result[7] + ", '" + result[8] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logSystemEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO system_event (time, server, transactionNum, command, userid, stocksymbol, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "', '" + result[7] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logAccountTransactionEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO account_transaction (time, server, transactionNum, action, userid, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logErrorEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO error_event (time, server, transactionNum, command, userid, stocksymbol, funds, errorMessage) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "' , '" + result[7] + "', '" + result[8] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logDebugEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO debug_event (time, server, transactionNum, command, userid, stocksymbol, funds, debugMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4] + "', " + result[5] + "', " + result[6] + "', " + result[7] + "', " + result[8] + "')").Exec(); err != nil {
		panic(err)
	}

}

func dumpUser(userId string, filename string) {

	var transactionTime string
	var server string
	var transactionNum string
	var command string
	var stockSymbol string
	var funds string
	filename = strings.TrimSpace(strings.TrimSuffix(filename, "\n"))
	filename = addTimestampToFilename(filename)

	log.Printf("Dumping user %s to file: %s\n", userId, filename)

	//check if user exists
	var count int
	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
		panic(err)
	}

	if count == 0 {
		log.Printf("User %s doesn't exist\n", userId)
		return
	}

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := etree.NewElement("log")
	doc.SetRoot(root)

	iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, stocksymbol, funds FROM usercommands WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&transactionTime, &server, &transactionNum, &command, &stockSymbol, &funds) {
		user_command(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
	}

	if err := iter.Close(); err != nil {
		panic(err)
	}

	doc.Indent(2)
	doc.WriteToFile(filename)

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
	var action string
	//var count int
	var errorMessage string
	var debugMessage string

	var servers string

	server_iter_user := sessionGlobal.Query("select distinct server from usercommands").Iter()
	for server_iter_user.Scan(&servers){
		
		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands where server ='" + servers + "'").PageSize(5000).Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			user_command(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}
	}

	//check if user commands
	/*
		if err := sessionGlobal.Query("SELECT count(*) FROM usercommands").Scan(&count); err != nil {
			panic(err)
		}
	

	count = 1
	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands").PageSize(5000).Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			user_command(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	*/
	//check if quote server events
	/*
		if err := sessionGlobal.Query("SELECT count(*) FROM quote_server").Scan(&count); err != nil {
			panic(err)
		}
	*/
	server_iter_quote := sessionGlobal.Query("select distinct server from quote_server").Iter()
	for server_iter_quote.Scan(&servers){
	//count = 1
	//if count != 0 {
		iter := sessionGlobal.Query("SELECT time, server, transactionNum, quoteservertime , userid, stocksymbol, price, cryptokey FROM quote_server where server ='" + servers + "'").PageSize(5000).Iter()
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
	*/

	server_iter_account_transaction := sessionGlobal.Query("select distinct server from account_transaction").Iter()

	for server_iter_account_transaction.Scan(&servers){
		iter := sessionGlobal.Query("SELECT time, server, transactionNum, action, userid, funds FROM account_transaction where server='" + servers + "'").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &action, &userId, &funds) {
			account_transaction(doc, transactionTime, server, transactionNum, action, userId, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if system event
	/*
	if err := sessionGlobal.Query("SELECT count(*) FROM system_event").Scan(&count); err != nil {
		panic(err)
	}
	*/
	server_iter_system := sessionGlobal.Query("select distinct server from system_event").Iter()

	for server_iter_system.Scan(&servers){
		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM system_event where server='" + servers + "'").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			system_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if error event
	/*
	if err := sessionGlobal.Query("SELECT count(*) FROM error_event").Scan(&count); err != nil {
		panic(err)
	}
	*/

	server_iter_error := sessionGlobal.Query("select distinct server from error_event").Iter()

	for server_iter_error.Scan(&servers){

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, errorMessage FROM error_event where server='" + servers + "'").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &errorMessage) {
			error_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds, errorMessage)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if debug event
	/*
	if err := sessionGlobal.Query("SELECT count(*) FROM debug_event").Scan(&count); err != nil {
		panic(err)
	}
	*/
	server_iter_debug := sessionGlobal.Query("select distinct server from debug_event").Iter()

	for server_iter_debug.Scan(&servers){
		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, debugMessage FROM debug_event where server='" + servers + "'").Iter()
		for iter.Scan(&transactionTime, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &debugMessage) {
			debug_event(doc, transactionTime, server, transactionNum, command, userId, stockSymbol, funds, debugMessage)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}

	doc.Indent(2)
	doc.WriteToFile(filename)
	fmt.Println("Done!")
}

func addTimestampToFilename(f string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", f, t.Format("2006-01-02_15-04-05"))
}
