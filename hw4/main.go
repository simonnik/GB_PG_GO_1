package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			err := getActivePosts(ctx, dbpool, 5)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)
			err = getPostComments(ctx, dbpool, 3)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)
			err = getUserById(ctx, dbpool, 3)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)
			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func getActivePosts(ctx context.Context, dbpool *pgxpool.Pool, limit int) error {
	const sql = `
SELECT * FROM posts
LEFT JOIN posts_images pi ON posts.id = pi.post_id
WHERE posts.is_active = true
LIMIT $1;
`
	rows, err := dbpool.Query(ctx, sql, limit)
	if err != nil {
		return fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	return nil
}

func getPostComments(ctx context.Context, dbpool *pgxpool.Pool, postId int) error {
	const sql = `SELECT * FROM posts_comments WHERE post_id = $1;`
	rows, err := dbpool.Query(ctx, sql, postId)
	if err != nil {
		return fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	return nil
}

func getUserById(ctx context.Context, dbpool *pgxpool.Pool, userId int) error {
	const sql = `SELECT * FROM users WHERE id = $1;`
	rows, err := dbpool.Query(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	return nil
}

func main() {
	ctx := context.Background()

	c := ReadConfig()
	pool, err := createPGXPool(int32(c.PoolMaxConns), int32(c.PoolMinConns))
	if err != nil {
		log.Fatalf("pool creation failed: %v", err)
	}
	defer pool.Close()

	duration := 10 * time.Second
	threads := 1000
	fmt.Println("start attack")
	res := attack(ctx, duration, threads, pool)
	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
}

func createPGXPool(maxConns int32, minConns int32) (*pgxpool.Pool, error) {
	cfg, err := getPoolConfig(maxConns, minConns)

	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	return pool, nil
}
func getPoolConfig(maxConns int32, minConns int32) (*pgxpool.Config, error) {
	url := "postgres://postgres:@localhost:54320/instabank"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pool config from connection string %s: %w", url, err)
	}

	cfg.MaxConns = maxConns
	cfg.MinConns = minConns

	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute

	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour

	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute

	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second

	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	return cfg, nil
}
