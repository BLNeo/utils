package template

import (
	"errors"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Redis struct {
	Addresses          []string `toml:"addresses"`
	MasterName         string   `toml:"master_name"`
	DialConnectTimeout string   `toml:"dial_connect_timeout"`
	DialReadTimeout    string   `toml:"dial_read_timeout"`
	DialWriteTimeout   string   `toml:"dial_write_timeout"`
	MaxIdle            string   `toml:"max_idle"`
	MaxActive          string   `toml:"max_active"`
	IdleTimeout        string   `toml:"idle_timeout"`
}

type RedisPool struct {
	Master *redis.Pool
	Slave  *redis.Pool
}

func (o *Redis) Engin(dbNumber int) (*RedisPool, error) {
	master, err := o.newRedis(true, dbNumber)
	if err != nil {
		return nil, err
	}
	slaver, err := o.newRedis(false, dbNumber)
	if err != nil {
		return nil, err
	}
	return &RedisPool{
		Master: master,
		Slave:  slaver,
	}, nil
}

func (o *Redis) newRedis(ismaster bool, dbNumber int) (*redis.Pool, error) {
	Log.Info("Redis Connection : "+strings.Join(o.Addresses, ",")+strconv.Itoa(dbNumber),
		zap.String("middleware", "Redis"))
	var newErr error
	sntnl := &sentinel.Sentinel{
		Addrs:      o.Addresses,
		MasterName: o.MasterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				newErr = err
				return nil, err
			}
			return c, nil
		},
	}
	rdb := &redis.Pool{
		MaxIdle:     100,
		MaxActive:   100,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error
			var addr string
			if ismaster {
				addr, err = sntnl.MasterAddr()
			} else {
				adds, err := sntnl.SlaveAddrs()
				if err != nil {
					return nil, err
				}
				n := rand.Intn(len(adds))
				addr = adds[n]
			}
			if err != nil {
				return nil, err
			}
			rand.Seed(time.Now().Unix())
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				newErr = err
				return nil, err
			}
			if dbNumber != 0 {
				_, err = c.Do("select", dbNumber)
				if err != nil {
					newErr = err
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				newErr = errors.New("Role check failed")
				return newErr
			} else {
				return nil
			}
		}}
	if newErr != nil {
		return nil, newErr
	}
	return rdb, nil
}
