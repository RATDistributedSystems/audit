package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

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
		panic(fmt.Sprintf("problem creating session", err))
	}

	addr, protocol := configs.GetServerDetails("audit")
	l, err := net.Listen(protocol, addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	var wg sync.WaitGroup
	fmt.Println("Listening on " + addr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, &wg)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, wg *sync.WaitGroup) {
	// will listen for message to process ending in newline (\n)
	print("received request")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	//remove new line character and any spaces received
	message = strings.TrimSuffix(message, "\n")
	message = strings.TrimSpace(message)
	//try to split the message into tokens for processing
	result := strings.Split(string(message), ",")
	//if not enough arguments, or incorrect format
	//send NA and close connection

	//add user event to database
	//time, server, transactionNum, command, userid, funds
	
	if result[0] == "User" {
		wg.Add(1)
		logUserEvent(result)
		wg.Done()

	}
	//add quote event to database
	//time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey
	if result[0] == "Quote" {
		wg.Add(1)
		logQuoteEvent(result)
		wg.Done()

	}
	//add system event to database
	//
	if result[0] == "System" {
		wg.Add(1)
		logSystemEvent(result)
		wg.Done()

	}
	//add error event to db
	if result[0] == "Error" {
		wg.Add(1)
		logErrorEvent(result)
		wg.Done()
	}

	//add account to db
	//time, server, transactionNum, action, userid, funds
	if result[0] == "Account" {
		wg.Add(1)
		logAccountTransactionEvent(result)
		wg.Done()


	}
	fmt.Println(result[0])
	if result[0] == "DUMPLOG" {
		if len(result) == 3 {
			//DUMP with specific user
			wg.Wait()
			dumpUser(result[1], result[2])
		} else if len(result) == 2 {
			//DUMP everything
			wg.Wait()
			dump(result[1])
		}
	}
	// Close the connection when you're done with it.
	conn.Close()
}

func logUserEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stockSymbol, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "' , '" + result[7] + "')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func logQuoteEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO quote_server (time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "' , " + result[7] + ", '" + result[8] + "')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func logSystemEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO system_event (time, server, transactionNum, command, userid, stocksymbol, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "', '" + result[7] + "')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func logAccountTransactionEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO account_transaction (time, server, transactionNum, action, userid, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4] + "', '" + result[5] + "', '" + result[6] + "')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func logErrorEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO error_event (time, server, transactionNum, command, userid, stocksymbols, funds, errorMessage) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + "', '" + result[4] + "', '" + result[5] + "', '" + result[6] + "', '" + result[7] + "', '" + result[8] + "', '" + result[9] + "')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func logDebugEvent(result []string, wg *sync.WaitGroup) {

	if err := sessionGlobal.Query("INSERT INTO debug_event (time, server, transactionNum, command, userid, stocksymbols, funds, debugMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4] + "', " + result[5] + "', " + result[6] + "', " + result[7] + "', " + result[8] + "', " + result[9] + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func dumpUser(userId string, filename string, wg *sync.WaitGroup) {

	wg.Wait()
	var time string
	var server string
	var transactionNum string
	var command string
	var stockSymbol string
	var funds string

	fmt.Println("In dump user with " + userId + " and " + filename)

	//check if user exists
	var count int
	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if count == 0 {
		fmt.Sprintf("No such user exists")
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
		panic(fmt.Sprintf("problem creating session", err))
	}

	doc.Indent(2)
	doc.WriteToFile(filename)

}

func dump(filename string, wg *sync.WaitGroup) {

	
	wg.Wait()
	fmt.Println("Starting dump log")

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := etree.NewElement("log")
	doc.SetRoot(root)

	filename = strings.TrimSuffix(filename, "\n")
	filename = strings.TrimSpace(filename)

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
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
	//check if quote server events
	if err := sessionGlobal.Query("SELECT count(*) FROM quote_server").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, quoteservertime , userid, stocksymbol, price, cryptokey FROM quote_server ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &quoteservertime, &userId, &stockSymbol, &price, &cryptokey) {
			quote_server(doc, time, server, transactionNum, quoteservertime, userId, stockSymbol, price, cryptokey)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
	//check if account transaction events
	if err := sessionGlobal.Query("SELECT count(*) FROM account_transaction").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, action, userid, funds FROM account_transaction ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &action, &userId, &funds) {
			account_transaction(doc, time, server, transactionNum, action, userId, funds)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
	//check if system event
	if err := sessionGlobal.Query("SELECT count(*) FROM system_event").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM system_event ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
			system_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
	//check if error event
	if err := sessionGlobal.Query("SELECT count(*) FROM error_event").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, errorMessage FROM error_event ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &errorMessage) {
			error_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds, errorMessage)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
	//check if debug event
	if err := sessionGlobal.Query("SELECT count(*) FROM debug_event").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	if count != 0 {

		iter := sessionGlobal.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, debugMessage FROM debug_event ").Iter()
		for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &debugMessage) {
			debug_event(doc, time, server, transactionNum, command, userId, stockSymbol, funds, debugMessage)
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}

	doc.Indent(2)
	doc.WriteToFile(filename)
	fmt.Println("Finished writing to file")

}
