package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
    "math/rand"
    "time"
    "strconv"
    "github.com/gocql/gocql"
    "github.com/beevik/etree"
)

const (
    CONN_HOST = "192.168.3.102"
    CONN_PORT = "5555"
    CONN_TYPE = "tcp"
)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
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
    if result[0] == "User"{
      go logUserEvent(result)
    }
    //add quote event to database
    //time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey
    if result[0] == "Quote"{
      go logQuoteEvent(result)
    }
    //add system event to database
    //
    if result[0] == "System"{
      go logSystemEvent(result)
    }
    //add error event to db
    if result[0] == "Error"{
      
    }

    //add account to db
    //time, server, transactionNum, action, userid, funds
    if result[0] == "Account"{
      
    }

    if result[0] == "DUMPLOG"{
      if len(result) == 3{
        //DUMP with specific user
        go dumpUser(result[1],result[2])
      } else if len(result) == 2 {
        //DUMP everything
        go dump(result[1])
      }
    }
  // Close the connection when you're done with it.
  conn.Close()
}

func logUserEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stockSymbol, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4]+ "', '" + result[5]+ "', '" + result[6] + "' , '" + result[7] + "')").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

}

func logQuoteEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO quote_server (time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4]+ "', '" + result[5]+ "', '" + result[6] + "' , " + result[7] + ", '" + result[8] + "')").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}

func logSystemEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO system_event (time, server, transactionNum, command, userid, stocksymbol, funds) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6]+ "', " + result[7] + ")").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}

func logAccountTransactionEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO account_transaction (time, server, transactionNum, action, userid, funds) VALUES (" + result[1] + ", '" + result[2] + "', '" + result[3] + "', '" + result[4]+ "', '" + result[5]+ "', '" + result[6] + "')").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}

func logErrorEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO error_event (time, server, transactionNum, command, userid, stocksymbols, funds, errorMessage) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + "', '" + result[4]+ "', '" + result[5]+ "', '" + result[6]+ "', '" + result[7]+ "', '" + result[8]+ "', '" + result[9] + "')").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}

func logDebugEvent(result []string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO debug_event (time, server, transactionNum, command, userid, stocksymbols, funds, debugMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6]+ "', " + result[7]+ "', " + result[8]+ "', " + result[9] + ")").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}


func dumpUser(userId string, filename string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  var time string
  var server string
  var transactionNum string
  var command string
  var stockSymbol string
  var funds string

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  //check if user exists
  var count int
  if err := session.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if count == 0{
    fmt.Sprintf("No such user exists")
    return
  }

  doc := etree.NewDocument()
  doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
  root := etree.NewElement("log")
  doc.SetRoot(root) 

  iter := session.Query("SELECT time, server, transactionNum, command, stocksymbol, funds FROM usercommands WHERE userid='" + userId + "'").Iter()
  for iter.Scan(&time, &server, &transactionNum, &command, &stockSymbol, &funds) {
      user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
    }
  
  if err := iter.Close(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  doc.Indent(2)
  doc.WriteToFile(filename)
  
}

func dump(filename string){
  
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  doc := etree.NewDocument()
  doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
  root := etree.NewElement("log")
  doc.SetRoot(root) 

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
  if err := session.Query("SELECT count(*) FROM usercommands").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){

    iter := session.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM usercommands ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
        user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    }  

  }
  //check if quote server events
  if err := session.Query("SELECT count(*) FROM quote_server").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){

    iter := session.Query("SELECT time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey FROM quote_server ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &price, &stockSymbol, &userId, &quoteservertime, &cryptokey) {
        quote_server(doc, time, server, transactionNum, price, stockSymbol, userId, quoteservertime, cryptokey)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    } 
    
  }
  //check if account transaction events
  if err := session.Query("SELECT count(*) FROM account_transaction").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){

    iter := session.Query("SELECT time, server, transactionNum, action, userid, funds FROM account_transaction ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &action, &userId, &funds) {
        account_transaction(doc, time, server, transactionNum, action, userId, funds)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    }     

  }
  //check if system event
  if err := session.Query("SELECT count(*) FROM system_event").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){

    iter := session.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds FROM system_event ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds) {
        user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    }   
    
  }
  //check if error event
  if err := session.Query("SELECT count(*) FROM error_event").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){
    
    iter := session.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, errorMessage FROM error_event ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &errorMessage) {
        user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    }  

  }
  //check if debug event
  if err := session.Query("SELECT count(*) FROM debug_event").Scan(&count); err != nil{
      panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
  }

  if (count != 0){
    
    iter := session.Query("SELECT time, server, transactionNum, command, userid, stocksymbol, funds, debugMessage FROM debug_event ").Iter()
    for iter.Scan(&time, &server, &transactionNum, &command, &userId, &stockSymbol, &funds, &debugMessage) {
        user_command(doc, time, server, transactionNum, command, userId, stockSymbol, funds)
      }
    
    if err := iter.Close(); err != nil {
      panic(fmt.Sprintf("problem creating session", err))
    } 

  }

  doc.Indent(2)
  doc.WriteToFile(filename)
    
}
