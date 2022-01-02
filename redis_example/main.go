package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const (
	redisHost   = ""
	redisPasswd = ""
)

var (
	pool *redis.Pool
)

func InitCache() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			if redisPasswd != "" {
				if _, err := c.Do("AUTH", redisPasswd); err != nil {
					c.Close()
					return nil, errors.Wrap(err, "")
				}
			}
			return c, errors.Wrap(err, "")
		},
	}
}

func CacheConn() redis.Conn {
	return pool.Get()
}

func main() {
	InitCache()

	Set("", "")

	getVal, _ := Get("")
	fmt.Printf(getVal)

}

func Do(key string, args ...interface{}) (interface{}, error) {
	conn := CacheConn()
	defer conn.Close()
	return conn.Do(key, args...)
}

func Get(key string) (string, error) {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return r, nil
}

func Del(key string) error {
	conn := CacheConn()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func Set(key, val string, expire ...int) error {
	conn := CacheConn()
	defer conn.Close()
	err := conn.Send("SET", key, val)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if len(expire) > 0 {
		err = conn.Send("EXPIRE", key, expire[0])
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	err = conn.Flush()
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func Lindex(key string, index int) (string, error) {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.String(conn.Do("LINDEX", key, index))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return r, nil
}

func Len(key string) int {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		return 0
	}
	return r
}

func Exists(key string) bool {
	conn := CacheConn()
	defer conn.Close()
	r, _ := redis.Int(conn.Do("EXISTS", key))
	if r <= 0 {
		return false
	}
	return true
}

//------------------LIST----------------------
func LPush(key string, val ...string) error {
	conn := CacheConn()
	defer conn.Close()
	_, err := conn.Do("LPUSH", key, val)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func RPop(key string) (string, error) {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return r, nil
}

//------------------SET----------------------
func SAdd(key string, val string) error {
	conn := CacheConn()
	defer conn.Close()
	_, err := conn.Do("SADD", key, val)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// 查看成员数
func Scard(key string) int {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SCARD", key))
	if err != nil {
		return 0
	}
	return r
}

// 判断是否存在于当前集合中
func SISMember(key, member string) int {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		return 0
	}
	return r
}

// 获取成员列表
func SMembers(key string) ([]string, error) {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 删除一个成员
func SRem(key string, member string) error {
	conn := CacheConn()
	defer conn.Close()
	_, err := conn.Do("SREM", key, member)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

//------------------HASH----------------------
// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func HMSet(key string, val interface{}, expire ...int) error {
	conn := CacheConn()
	defer conn.Close()
	// _, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(val))
	conn.Send("MULTI")
	err := conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	if len(expire) > 0 {
		conn.Send("PEXPIRE")
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return errors.Wrap(err, "")
	}
	// conn.Flush()
	// if err != nil {
	// 	return errors.Wrap(err, "")
	// }
	// _, err = conn.Receive()
	return nil
}

// 将哈希表 key 中的字段 field 的值设为 value 。
func HSet(key, field string, val interface{}) error {
	conn := CacheConn()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, val)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func HGetAll(key string) (map[string]string, error) {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func HGet(key, field string) string {
	conn := CacheConn()
	defer conn.Close()
	r, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return ""
	}
	return r
}
