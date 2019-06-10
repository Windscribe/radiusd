package config

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

type Listener struct {
	Addr   string
	Secret string
	CIDR   []string
}

type Conf struct {
	Dsn           string
	Listen        map[string]Listener
	ControlListen string
}

var (
	C         *Conf
	Log       *log.Logger
	Debug     bool
	Verbose   bool
	Hostname  string
	DB        *sql.DB
	ErrNoRows = sql.ErrNoRows
	Stopping  bool
	Sock      []*net.UDPConn
)

func Init(path string) error {
	r, e := os.Open(path)
	if e != nil {
		return e
	}
	defer r.Close()

	C = new(Conf)
	if _, e := toml.DecodeReader(r, &C); e != nil {
		return fmt.Errorf("TOML: %s", e)
	}
	Hostname, e = os.Hostname()
	if e != nil {
		panic(e)
	}

	Log = log.New(os.Stdout, "radiusd ", log.LstdFlags)
	return dbInit("mysql", C.Dsn)
}

func dbInit(driver string, dsn string) error {
	var e error
	DB, e = sql.Open(driver, dsn)
	if e != nil {
		return e
	}
	return DB.Ping()
}

func DbClose() error {
	return DB.Close()
}
