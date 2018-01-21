package main

import (

	"github.com/beevik/etree"
)



func main(){
	
	
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)


	user_command(doc, "9433", "SN6", "89", "BUY", "user1", "300.00")
	quote_server(doc, "9838", "N6", "90", "08943", "user2", "XYZ", "99.66", "kjhdf832ihfkl8fj")
	system_event(doc, "474784", "Z89", "92", "QUOTE", "system1", "BNM", "66.88")
	error_event(doc, "3904", "S8", "93", "BUY", "user4", "GTY", "999.00", "Insufficient funds")
	account_transaction(doc, "040302", "P9", "91", "SELL", "user1", "299.99")



	doc.Indent(2)
	doc.WriteToFile("log_sample.xml")


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
	timestamp.CreateCharData(ts)
	server := accountTransaction.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := accountTransaction.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	action := accountTransaction.CreateElement("action")
	action.CreateCharData(act)
	username := accountTransaction.CreateElement("username")
	username.CreateCharData(uname)
	funds := accountTransaction.CreateElement("funds")
	funds.CreateCharData(fund)
}
//system events. current user commands, inter"server" communication, execution trigger
func system_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string){
	systemEvent := doc.CreateElement("systemEvent")
	timestamp := systemEvent.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := systemEvent.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := systemEvent.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := systemEvent.CreateElement("command")
	command.CreateCharData(cmd)
	username := systemEvent.CreateElement("username")
	username.CreateCharData(uname)
	stockSymbol := systemEvent.CreateElement("stockSymbol")
	stockSymbol.CreateCharData(stocksym)
	funds := systemEvent.CreateElement("funds")
	funds.CreateCharData(fund)
}
//error messages, user commands, optional error message
func error_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string, err string){
	errorEvent := doc.CreateElement("errorEvent")
	timestamp := errorEvent.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := errorEvent.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := errorEvent.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := errorEvent.CreateElement("command")
	command.CreateCharData(cmd)
	username := errorEvent.CreateElement("username")
	username.CreateCharData(uname)
	stockSymbol := errorEvent.CreateElement("stockSymbol")
	stockSymbol.CreateCharData(stocksym)
	funds := errorEvent.CreateElement("funds")
	funds.CreateCharData(fund)
	errorMessage := errorEvent.CreateElement("errorMessage")
	errorMessage.CreateCharData(err)
}

//debugging messages, user commands, optional debug message


