package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	gotime "time"

	"github.com/RATDistributedSystems/utilities"
	"github.com/beevik/etree"
	"github.com/gocql/gocql"
)

var sessionGlobal *gocql.Session
var configs = utilities.GetConfigurationFile("config.json")

func main() {
	//establish global connection param

	cluster := gocql.NewCluster(configs.GetValue("cassandra_ip"))
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

	addr, protocol := configs.GetServerDetails("audit")
	l, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + addr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	message = strings.TrimSpace(strings.TrimSuffix(message, "\n"))
	log.Printf("Received request: %s\n", message)
	result := strings.Split(string(message), ",")
	defer conn.Close()

	if result[0] == "User" {
		logUserEvent(result)
	}

	if result[0] == "Quote" {
		logQuoteEvent(result)
	}

	if result[0] == "System" {
		logSystemEvent(result)
	}

	if result[0] == "Error" {
		logErrorEvent(result)
	}

	if result[0] == "Account" {
		logAccountTransactionEvent(result)
	}

	if result[0] == "DUMPLOG" {
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

	if err := sessionGlobal.Query("INSERT INTO error_event (time, server, transactionNum, command, userid, stocksymbols, funds, errorMessage) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + "', '" + result[4] + "', '" + result[5] + "', '" + result[6] + "', '" + result[7] + "', '" + result[8] + "', '" + result[9] + "')").Exec(); err != nil {
		panic(err)
	}

}

func logDebugEvent(result []string) {

	if err := sessionGlobal.Query("INSERT INTO debug_event (time, server, transactionNum, command, userid, stocksymbols, funds, debugMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4] + "', " + result[5] + "', " + result[6] + "', " + result[7] + "', " + result[8] + "', " + result[9] + ")").Exec(); err != nil {
		panic(err)
	}

}

func dumpUser(userId string, filename string) {

	var time string
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
	for iter.Scan(&time, &server, &transactionNum, &command, &stockSymbol, &funds) {
		user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
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

	var time string
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
	var count int
	var errorMessage string
	var debugMessage string

	//check if user commands
	if err := sessionGlobal.Query("SELECT count(*) FROM usercommands").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if quote server events
	if err := sessionGlobal.Query("SELECT count(*) FROM quote_server").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, quoteservertime , userid, stocksymbol, price, cryptokey FROM quote_server ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &quoteservertime, &userId, &stockSymbol, &price, &cryptokey) {
			quote_server(doc, time, server, transactionNum, quoteservertime, userId, stockSymbol, price, cryptokey)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}
	//check if account transaction events
	if err := sessionGlobal.Query("SELECT count(*) FROM account_transaction").Scan(&count); err != nil {
		panic(err)
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, action, userid, funds FROM account_transaction ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &action, &userId, &funds) {
			account_transaction(doc, time, server, transactionNum, action, userId, funds)
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
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			system_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
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
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &errorMessage) {
			error_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds, errorMessage)
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
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &debugMessage) {
			debug_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds, debugMessage)
		}

		if err := iter.Close(); err != nil {
			panic(err)
		}

	}

	doc.Indent(2)
	doc.WriteToFile(filename)
}

func addTimestampToFilename(f string) string {
	t := gotime.Now()
	return fmt.Sprintf("%s-%s", t.Format("2006-01-02_15-04-05"), f)
}
