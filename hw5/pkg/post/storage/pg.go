package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/simonnik/GB_PG_GO_1/hw5/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	db    *pgxpool.Pool
	dbMux = &sync.Mutex{}
)

type DB interface {
	Close()
	GetPosts() ([]*Post, error)
}

type conn struct {
	db *pgxpool.Pool
}

type Image struct {
	PostId uint   `json:"-" db:"post_id"`
	Path   string `json:"path"`
}
type Post struct {
	Id          uint        `json:"id"`
	CreatedAt   interface{} `json:"created_at"`
	Description string      `json:"description"`
	Images      []*Image    `json:"images"`
}

func NewDB(cfg *config.Config) (DB, error) {
	pool, err := getConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get a connection pool: %w", err)
	}
	return &conn{
		db: pool,
	}, nil
}

func (c *conn) Close() {
	c.db.Close()
}

func (c *conn) GetPosts() ([]*Post, error) {
	posts := make([]*Post, 0)
	sql := `SELECT id, created_at, description FROM posts
			WHERE is_active = true
			LIMIT 50;`
	if err := pgxscan.Select(context.Background(), c.db, &posts, sql); err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}

	var postIds string
	for i, post := range posts {
		if i > 0 {
			postIds += ", "
		}
		postIds = strconv.Itoa(int(post.Id))
	}

	imgSql := `SELECT post_id, path FROM posts_images WHERE post_id IN ($1)`
	imgs := make([]*Image, 0)

	if err := pgxscan.Select(context.Background(), c.db, &imgs, imgSql, postIds); err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	images := map[uint][]*Image{}
	for _, img := range imgs {
		images[img.PostId] = append(images[img.PostId], img)
	}

	for _, post := range posts {
		img, ok := images[post.Id]
		if !ok {
			continue
		}

		post.Images = img
	}
	return posts, nil
}

func getConn(cfg *config.Config) (*pgxpool.Pool, error) {
	dbMux.Lock()
	defer dbMux.Unlock()
	if db != nil {
		return db, nil
	}

	var err error
	db, err = initPGXPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a PGX pool: %w", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return db, nil
}

func initPGXPool(c *config.Config) (*pgxpool.Pool, error) {
	connStr, err := composeConnectionString(c)
	if err != nil {
		return nil, fmt.Errorf("failed to compose the connection string: %w", err)
	}
	cfg, err := getPGXPoolConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the PGX pool config: %w", err)
	}
	db, err = pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the postgres DB using a PGX connection pool: %w", err)
	}
	return db, nil
}

func getPGXPoolConfig(connStr string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create the PGX pool config from connection string: %w", err)
	}
	cfg.ConnConfig.ConnectTimeout = time.Second * 1
	cfg.ConnConfig.Logger = logrusadapter.NewLogger(
		&logrus.Logger{
			Out:          os.Stdout,
			Formatter:    new(logrus.JSONFormatter),
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.InfoLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		})
	return cfg, nil
}

func composeConnectionString(c *config.Config) (string, error) {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		url.QueryEscape(c.DB.User),
		url.QueryEscape(c.DB.Password),
		url.QueryEscape(c.DB.Host),
		url.QueryEscape(c.DB.Port),
		url.QueryEscape(c.DB.Name),
	), nil
}
