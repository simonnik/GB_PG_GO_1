package main

import "flag"

type Config struct {
	PoolMaxConns    int
	PoolMinConns    int
	GoroutinesCount int
}

func ReadConfig() *Config {
	c := &Config{}
	flag.IntVar(&c.PoolMaxConns, "max-conns", 8, "connection pool MaxConnections param")
	flag.IntVar(&c.PoolMaxConns, "min-conns", 8, "connection pool MinConnections param")
	flag.IntVar(&c.GoroutinesCount, "goroutines", 50, "number of goroutines to run")
	flag.Parse()
	return c
}
