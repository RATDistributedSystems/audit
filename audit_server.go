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

    }
    //add system event to database
    //
    if result[0] == "System"{

    }
    //add error event to db
    if result[0] == "Error"{
      
    }

    //add account to db
    //time, server, transactionNum, action, userid, funds
    if result[0] == "Account"{
      
    }

    if result[0] == "DUMP"{
      if len(result) == 2{
        //DUMP with specific user
        go dumpUser(result[1])
      } else {
        //DUMP everything
      }
    }

    if len(result) != 2 || len(result[0]) != 3 {
      conn.Write([]byte("NA"))
      conn.Close()
    //correct input
    } else {
      username := result[1]
      username = strings.TrimSpace(username)
      stock_sym := result[0]
      stock_sym = strings.TrimSpace(stock_sym)
      r := rand.New(rand.NewSource(time.Now().UnixNano()))
      //generate random price
      rand_price := r.Float64() * float64(rand.Intn(1000))
      stock_price := strconv.FormatFloat(rand_price, 'f', 2, 64)
      //get curret timestamp in UTC
      t := time.Now().UTC().UnixNano()
      time := strconv.Itoa(int(t))
      //generate cryptokey
      crypto := "7777777777"
      /*
      print(stock_price + ",")
      print(result[0] +",")
      print(result[1] + ",")
      print(time + ",")
      print(crypto)
      */
      fmt.Fprintf(conn, stock_price + "," + stock_sym + "," + username + "," + time + "," + crypto)

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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, funds) VALUES (" + result[1] + ", '" + result[2] + "', " + result[3] + ", '" + result[4]+ "', '" + result[5]+ "', '" + result[6] + "')").Exec(); err != nil {
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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, price, stocksymbol, userid, quoteservertime, cryptokey) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6] + result[7] + "', " + result[8] + ")").Exec(); err != nil {
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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stocksymbol, funds) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6]+ "', " + result[7] + ")").Exec(); err != nil {
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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, action, userid, funds) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6] + ")").Exec(); err != nil {
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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stocksymbols, funds, errorMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6]+ "', " + result[7]+ "', " + result[8]+ "', " + result[9] + ")").Exec(); err != nil {
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

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, stocksymbols, funds, debugMessage) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6]+ "', " + result[7]+ "', " + result[8]+ "', " + result[9] + ")").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  
}


func dumpUser(userId string){
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  var time string
  var server string
  var transactionNum string
  var command string
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

  iter := session.Query("SELECT time, server, transactionNum, command, funds FROM usercommands WHERE userid='" + userId + "'").Iter()
  for iter.Scan(&time, &server, &transactionNum, &command, &funds) {
      user_command(doc, time, server, transactionNum, command, userId, funds)
    }
  
  if err := iter.Close(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  doc.Indent(2)
  doc.WriteToFile(userId + "_log.xml")
  
}

func dump(){
  /*
  cluster := gocql.NewCluster("192.168.3.103")
  cluster.Keyspace = "userdb"
  cluster.ProtoVersion = 4
  session, err := cluster.CreateSession()

  if err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }

  if err := session.Query("INSERT INTO usercommands (time, server, transactionNum, command, userid, funds) VALUES ('" + result[1] + "', " + result[2] + "', " + result[3] + "', " + result[4]+ "', " + result[5]+ "', " + result[6] + ")").Exec(); err != nil {
    panic(fmt.Sprintf("problem creating session", err))
  }
  */
    
}
