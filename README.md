Justaway Ex
=============

Justaway Extension API Server

## Install

### crawler

### web

```bash
cd web
export JUSTAWAY_EX_CONSUMER_KEY=''
export JUSTAWAY_EX_CONSUMER_SECRET=''
export JUSTAWAY_EX_DB_SOURCE='justaway@tcp(192.168.0.10:3306)/justaway'
export JUSTAWAY_EX_CALLBACK='http://127.0.0.1:8002/signin/callback'
go run main.go
```

```bash
cd crawler
export JUSTAWAY_EX_CONSUMER_KEY=''
export JUSTAWAY_EX_CONSUMER_SECRET=''
export JUSTAWAY_EX_DB_SOURCE='justaway@tcp(192.168.0.10:3306)/justaway'
export JUSTAWAY_EX_CRAWLER_ID='1'
go run *.go

mysql -h 192.168.0.10 -u justaway justaway
> echo "INSERT INTO crawler(url, created_at, updated_at) VALUES('http://127.0.0.1:8001/', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW()));"
```
