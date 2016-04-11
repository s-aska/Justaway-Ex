Justaway Ex
=============

Justaway Extension API Server

## Install

### crawler

```bash
cd crawler
export CONSUMER_KEY=''
export CONSUMER_SECRET=''
export JUSTAWAY_EX_DB_SOURCE='justaway@tcp(192.168.0.10:3306)/justaway'
export JUSTAWAY_EX_CRAWLER_ID=1
go run *.go
```

### web

```bash
cd web
export CONSUMER_KEY=''
export CONSUMER_SECRET=''
export JUSTAWAY_EX_DB_SOURCE='justaway@tcp(192.168.0.10:3306)/justaway'
export JUSTAWAY_EX_CALLBACK='http://127.0.0.1:8002/signin/callback'
go run *.go
```

