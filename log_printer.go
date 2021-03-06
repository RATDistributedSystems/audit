package main

import (
	"github.com/beevik/etree"
)

//receive an user command
func user_command(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stockSym string, fund string) {
	root := doc.Root()
	userCommand := root.CreateElement("userCommand")
	timestamp := userCommand.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := userCommand.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := userCommand.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := userCommand.CreateElement("command")
	command.CreateCharData(cmd)
	if (uname != "") && (uname != "-1") {
		username := userCommand.CreateElement("username")
		username.CreateCharData(uname)
	}
	if stockSym != "" {
		stockSymbol := userCommand.CreateElement("stockSymbol")
		stockSymbol.CreateCharData(stockSym)
	}
	if fund != "" {
		funds := userCommand.CreateElement("funds")
		funds.CreateCharData(fund)
	}
}

//hit to quote server
func quote_server(doc *etree.Document, ts string, srvr string, trsNum string, quoteTime string, uname string, stocksym string, pric string, crypto string) {
	root := doc.Root()
	quoteServer := root.CreateElement("quoteServer")
	timestamp := quoteServer.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := quoteServer.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := quoteServer.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	quoteServerTime := quoteServer.CreateElement("quoteServerTime")
	quoteServerTime.CreateCharData(quoteTime)
	username := quoteServer.CreateElement("username")
	username.CreateCharData(uname)
	stockSymbol := quoteServer.CreateElement("stockSymbol")
	stockSymbol.CreateCharData(stocksym)
	price := quoteServer.CreateElement("price")
	price.CreateCharData(pric)
	cryptokey := quoteServer.CreateElement("cryptokey")
	cryptokey.CreateCharData(crypto)
}

//user account touch. Add remove
func account_transaction(doc *etree.Document, ts string, srvr string, trsNum string, act string, uname string, fund string) {
	root := doc.Root()
	accountTransaction := root.CreateElement("accountTransaction")
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
func system_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string) {
	root := doc.Root()
	systemEvent := root.CreateElement("systemEvent")
	timestamp := systemEvent.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := systemEvent.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := systemEvent.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := systemEvent.CreateElement("command")
	command.CreateCharData(cmd)
	if uname != "" {
		username := systemEvent.CreateElement("username")
		username.CreateCharData(uname)
	}
	if stocksym != "" {
		stockSymbol := systemEvent.CreateElement("stockSymbol")
		stockSymbol.CreateCharData(stocksym)
	}
	if fund != "" {
		funds := systemEvent.CreateElement("funds")
		funds.CreateCharData(fund)
	}

}

//error messages, user commands, optional error message
func error_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string, err string) {
	root := doc.Root()
	errorEvent := root.CreateElement("errorEvent")
	timestamp := errorEvent.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := errorEvent.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := errorEvent.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := errorEvent.CreateElement("command")
	command.CreateCharData(cmd)
	if uname != "" {
		username := errorEvent.CreateElement("username")
		username.CreateCharData(uname)
	}

	if stocksym != "" {
		stockSymbol := errorEvent.CreateElement("stockSymbol")
		stockSymbol.CreateCharData(stocksym)
	}

	if fund != "" {
		funds := errorEvent.CreateElement("funds")
		funds.CreateCharData(fund)
	}

	if err != "" {
		errorMessage := errorEvent.CreateElement("errorMessage")
		errorMessage.CreateCharData(err)
	}

}

//debugging messages, user commands, optional debug message
func debug_event(doc *etree.Document, ts string, srvr string, trsNum string, cmd string, uname string, stocksym string, fund string, debug string) {
	root := doc.Root()
	debugEvent := root.CreateElement("debugEvent")
	timestamp := debugEvent.CreateElement("timestamp")
	timestamp.CreateCharData(ts)
	server := debugEvent.CreateElement("server")
	server.CreateCharData(srvr)
	transactionNum := debugEvent.CreateElement("transactionNum")
	transactionNum.CreateCharData(trsNum)
	command := debugEvent.CreateElement("command")
	command.CreateCharData(cmd)
	if uname != "" {
		username := debugEvent.CreateElement("username")
		username.CreateCharData(uname)
	}

	if stocksym != "" {
		stockSymbol := debugEvent.CreateElement("stockSymbol")
		stockSymbol.CreateCharData(stocksym)
	}

	if fund != "" {
		funds := debugEvent.CreateElement("funds")
		funds.CreateCharData(fund)
	}

	if debug != "" {
		errorMessage := debugEvent.CreateElement("errorMessage")
		errorMessage.CreateCharData(debug)
	}
}
