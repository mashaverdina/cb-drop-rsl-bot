package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"golang.yandex/hasql"
	"golang.yandex/hasql/checkers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PGClient struct {
	cluster *hasql.Cluster
}

var (
	ErrMasterIsUnavailable = errors.New("postgres: master is unavailable")
	ErrNodeIsUnavailable   = errors.New("postgres: node is unavailable")
)

func NewPGClient(hosts []string, port int, user string, password string, db string) (*PGClient, error) {
	var nodes []hasql.Node
	for _, host := range hosts {
		connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db)
		connConfig, err := pgx.ParseConfig(connString)
		if err != nil {
			return nil, err
		}

		connConfig.BuildStatementCache = nil
		connConfig.PreferSimpleProtocol = true

		db := stdlib.OpenDB(*connConfig)
		nodes = append(nodes, hasql.NewNode(host, db))
	}

	opts := []hasql.ClusterOption{
		hasql.WithUpdateInterval(2 * time.Second),
		hasql.WithNodePicker(hasql.PickNodeRoundRobin()),
		hasql.WithTracer(newHASQLTracer()),
	}
	c, err := hasql.NewCluster(nodes, checkers.PostgreSQL, opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := c.WaitForPrimary(ctx); err != nil {
		return nil, err
	}

	return &PGClient{c}, nil
}

func newHASQLTracer() hasql.Tracer {
	return hasql.Tracer{
		UpdateNodes: func() {
			// log.Println("update nodes")
		},
		UpdatedNodes: func(nodes hasql.AliveNodes) {
			// log.Printf("nodes were updates: %+v\n", nodes)
		},
		NodeDead: func(node hasql.Node, err error) {
			log.Printf("node %q is dead: %v\n", node, err)
		},
		NodeAlive: func(node hasql.Node) {
			// log.Printf("node %q is alive\n", node)
		},
		NotifiedWaiters: func() {
			// log.Println("notified all waiters")
		},
	}
}

func (c *PGClient) GetMaster() (*gorm.DB, error) {
	if node := c.cluster.Primary(); node != nil {
		return gormWrapper(node)
	}
	return nil, ErrMasterIsUnavailable
}

func (c *PGClient) ExecuteInTransaction(nodeState hasql.NodeStateCriteria, processFunc func(*gorm.DB) error) error {
	if node := c.cluster.Node(nodeState); node != nil {
		db, err := gormWrapper(node)
		if err != nil {
			return err
		}
		return db.Transaction(processFunc)
	}

	return ErrNodeIsUnavailable
}

func gormWrapper(node hasql.Node) (*gorm.DB, error) {
	return gorm.Open(
		postgres.New(
			postgres.Config{
				Conn: node.DB(),
			},
		), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{SlowThreshold: 2 * time.Second, LogLevel: logger.Info},
			),
		},
	)
}
