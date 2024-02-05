package main

import { 
	"fmt"
	"sync"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
}

func newSqlConn() *sql.DB {
	// build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	_db, err :=sql.Open("mysql", dsn)
	
	if err != nil {
		panic(err)
	} 
	return _db
}

type conn struct{
	db *sql.DB
}

type cpool struct{
	mu *sync.Mutex
	channel chan interface()
	conns []*conn
	maxConn int
}

func newSqlConnPool(maxConn int) (*cpool, error){
	var mu = sync.Mutex{}
	pool := &cpool{
		mu: &mu,
		conns: make([]*conn, 0, maxConn)
		maxConn: maxConn,
		channel: make(chan interface{}, maxConn),
	}

	for i:=0; i < maxConn; i++ {
		pool.conns = append(pool.conns, &conn{newSqlConnPool()})
		pool.channel <- nil
	}
	
	return pool, nil
}

func (pool *cpool) Close(){
	close(pool.channel)
	for i:=range pool.conns {
		pool.conns[i].db.Close()
	}
}



func main() {
	fmt.Println("Hello!")
}
