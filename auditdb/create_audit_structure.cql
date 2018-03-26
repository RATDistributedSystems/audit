CREATE KEYSPACE userdb WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND durable_writes = true;

CREATE TABLE userdb.usercommands (
    time bigint,
    server text,
    transactionNum int,
    command text,
    userid text,
    stockSymbol text,
    funds text,
    PRIMARY KEY((server), transactionNum)
);

CREATE TABLE userdb.quote_server (
    time bigint,
    server text,
    transactionNum int,
    quoteservertime bigint,
    userid text,
    stocksymbol text,
    price text,
    cryptokey text,
    PRIMARY KEY((userid), time)
);

CREATE TABLE userdb.account_transaction(
    time bigint,
    server text,
    transactionNum int,
    action text,
    userid text,
    funds text,
    PRIMARY KEY((userid), time)
);

CREATE TABLE userdb.system_event(
    time bigint,
    server text,
    transactionNum int,
    command text,
    userid text,
    stocksymbol text,
    funds text,
    PRIMARY KEY((userid), time)
);

CREATE TABLE userdb.error_event(
    time bigint,
    server text,
    transactionNum int,
    command text,
    userid text,
    stocksymbol text,
    funds text,
    errorMessage text,
    PRIMARY KEY((userid), time)
);


CREATE TABLE userdb.debug_event(
    time bigint,
    server text,
    transactionNum int,
    command text,
    userid text,
    stocksymbol text,
    filename text,
    funds text,
    debugMessage text,
    PRIMARY KEY((userid), time)
);
