package postgresql

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	//	pgx functions :)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	// SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	// CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	// BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	// Ping(ctx context.Context) error
}

// structure used to create new psql client
type ClientConfig struct {
	Login    string
	Password string
	Host     string
	Port     int
	DBName   string

	PoolConfig
}

type PoolConfig struct {
	connStirng string

	MaxConns                 int32
	MaxConnIdleMinutes       int
	MaxConnLifetimeMinutes   int
	HealthCheckPeriodMinutes int
}

var (
	pool *pgxpool.Pool
	once sync.Once
)

// Constructor for postgres client
func NewClient(ctx context.Context, req ClientConfig) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		req.connStirng = fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s",
			req.Login,
			req.Password,
			req.Host,
			req.Port,
			req.DBName,
		)
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		pool, err = getConnectionPool(ctx, req.PoolConfig)
		if err != nil {
			log.Fatal("Error while getting pool connection: ", err)
		}
	})

	return pool, err
}

func getConnectionPool(ctx context.Context, configRequest PoolConfig) (pool *pgxpool.Pool, err error) {
	cfg, err := pgxpool.ParseConfig(configRequest.connStirng)
	if err != nil {
		println("Error while parsing config from connection string")
		return
	}

	cfg.MaxConns = configRequest.MaxConns
	cfg.MaxConnIdleTime = time.Minute * time.Duration(configRequest.MaxConnIdleMinutes)
	cfg.MaxConnLifetime = time.Minute * time.Duration(configRequest.MaxConnLifetimeMinutes)
	cfg.HealthCheckPeriod = time.Minute * time.Duration(configRequest.HealthCheckPeriodMinutes)
	pool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		println("Error while connecting to postgres by config")
		return
	}

	return
}

func NamedQuery(c Client, ctx context.Context, sql string, args map[string]interface{}) (pgx.Rows, error) {
	query, arguments := replaceQueryParams(sql, args)
	return c.Query(ctx, query, arguments...)
}

func replaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" && strings.Contains(namedQuery, ":"+k) {
			namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))
			args = append(args, v)
			i++
		}
	}

	return namedQuery, args
}
