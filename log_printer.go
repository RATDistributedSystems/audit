package main

import (

	"github.com/beevik/etree"
)



func main(){
	
	
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)


	user_command(doc, "9433", "SN6", "89", "BUY", "8947fjkfj", "300.00")

	doc.Indent(2)
	doc.WriteToFile("name")


}


//receive an user command
func user_command(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, fund string){
	userCommand := doc.CreateElement("userCommand")
	timestamp := userCommand.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := userCommand.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := userCommand.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := userCommand.CreateElement("command")
	command.CreateCharData(cmd)
	username := userCommand.CreateElement("username")
	username.CreateCharData(uname)
	funds := userCommand.CreateElement("funds")
	funds.CreateCharData(fund)


}

//hit to quote server
func quote_server(doc *etree.Document, ts string, srvr string, trsNum string, quoteTime string, uname string, stocksym string, pric string, crypto string){
	quoteServer := doc.CreateElement("quoteServer")
	timestamp := quoteServer.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := quoteServer.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := quoteServer.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	quoteServerTime := quoteServer.CreateElement("quoteTime")
	quoteServerTime.CreateCharData(quoteTime)
	username := quoteServer.CreateElement("username")
	username.CreateCharData(uname)
	stockSymbol := quoteServer.CreateElement("stockSymbol")
	stockSymbol.CreateCharData(stocksym)
	price	:= quoteServer.CreateElement("price")
	price.CreateCharData(pric)
	cryptokey := quoteServer.CreateElement("cryptokey")
	cryptokey.CreateCharData(crypto)

}

//user account touch. Add remove
func account_transaction(doc *etree.Document, ts string, srvr string, trsNum string, act string, uname string, fund string){
	accountTransaction := doc.CreateElement("accountTransaction")
	timestamp := accountTransaction.CreateElement("timestamp")
	server := accountTransaction.CreateElement("server")
	transactionNum := accountTransaction.CreateElement("transactionNum")
	action := accountTransaction.CreateElement("action")
	username := accountTransaction.CreateElement("username")
	funds := accountTransaction.CreateElement("funds")
}

//system events. current user commands, inter"server" communication, execution trigger
func system_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string){
	systemEvent := doc.CreateElement("systemEvent")
	timestamp := systemEvent.CreateElement("timestamp")
	server := systemEvent.CreateElement("server")
	transactionNum := systemEvent.CreateElement("transactionNum")
	command := systemEvent.CreateElement("command")
	username := systemEvent.CreateElement("username")
	stockSymbol := systemEvent.CreateElement("stockSymbol")
	funds := systemEvent.CreateElement("funds")
}

//error messages, user commands, optional error message
func error_event(doc *etree.Document){
	errorEvent := doc.CreateElement("errorEvent")
	errorEvent.CreateElement("timestamp")
	errorEvent.CreateElement("server")
	errorEvent.CreateElement("transactionNum")
	errorEvent.CreateElement("command")
	errorEvent.CreateElement("username")
	errorEvent.CreateElement("stockSymbol")
	errorEvent.CreateElement("funds")
	errorEvent.CreateElement("errorMessage")
}

//debugging messages, user commands, optional debug message


