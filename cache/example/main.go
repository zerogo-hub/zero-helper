package main

import (
	"errors"

	zerocache "github.com/zerogo-hub/zero-helper/cache"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

var (
	// RedisHost ...
	RedisHost = "127.0.0.1"
	// RedisPort ...
	RedisPort = 11234
	// RedisPassword ...
	RedisPassword = "uio876..."
)

var log = zerologger.NewSampleLogger()

func main() {

	c := zerocache.NewCache(
		zerocache.WithHost(RedisHost),
		zerocache.WithPort(RedisPort),
		zerocache.WithPassword(RedisPassword),
	)

	err := c.Open()
	if err != nil {
		log.Error("cache open failed: %s", err.Error())
		return
	}

	if err := testString(c); err != nil {
		log.Errorf("testString failed: %s", err.Error())
		return
	}

	if err := testHash(c); err != nil {
		log.Errorf("testHash failed: %s", err.Error())
		return
	}

	if err := testList(c); err != nil {
		log.Errorf("testList failed: %s", err.Error())
		return
	}

	if err := testSet(c); err != nil {
		log.Errorf("testSet failed: %s", err.Error())
		return
	}

	if err := testSortedSet(c); err != nil {
		log.Errorf("testSortedSet failed: %s", err.Error())
		return
	}

	log.Info("test cache success")
}

func testString(c zerocache.Cache) error {
	const (
		key   = "hello"
		value = "world"
	)

	if err := c.Set(key, value); err != nil {
		return err
	}

	bv, _ := c.GetBytes(key)
	if string(bv) != value {
		return errors.New("test GetBytes failed")
	}

	ttl, err := c.TTL(key)
	if err != nil {
		return err
	}
	if ttl != -1 {
		return errors.New("testString error 1")
	}

	v, err := c.Get(key)
	if err != nil {
		return err
	}
	if v != value {
		return errors.New("testString error 2")
	}

	if err := c.SetEx(key, value, "10"); err != nil {
		return err
	}

	ttl, err = c.TTL(key)
	if err != nil {
		return err
	}
	if ttl <= 0 || ttl > 10 {
		return errors.New("error 3")
	}

	if err := c.PSetEx(key, value, "10000"); err != nil {
		return err
	}

	ttl, err = c.PTTL(key)
	if err != nil {
		return err
	}
	if ttl <= 0 || ttl > 10000 {
		return errors.New("error 4")
	}

	if err := c.MSet("key-1", "value-1", "key-2", "value-2"); err != nil {
		return err
	}

	vs, err := c.MGet("key-1", "key-2", "key-3")
	if err != nil {
		return err
	}
	if len(vs) != 3 {
		return errors.New("testString error 5")
	}

	if err := c.Set(key, value); err != nil {
		return err
	}

	n, err := c.Strlen(key)
	if err != nil {
		return err
	}
	if n != len(key) {
		return errors.New("testString error 6")
	}
	n, err = c.Append(key, "haha")
	if err != nil {
		return err
	}
	n2, _ := c.Strlen(key)
	if n != n2 {
		return errors.New("testString error 7")
	}

	if err := c.Set(key, value); err != nil {
		return err
	}

	newValue := "gogogo"
	v, err = c.GetSet(key, newValue)
	if err != nil {
		return err
	}
	if v != value {
		return errors.New("testString error 8")
	}
	v, err = c.Get(key)
	if err != nil {
		return err
	}
	if v != newValue {
		return errors.New("testString error 9")
	}

	var (
		nKey   = "num"
		nValue = "3"
	)

	if err := c.Set(nKey, nValue); err != nil {
		return err
	}

	n64, err := c.Incr(nKey)
	if err != nil {
		return err
	}
	if n64 != 4 {
		return errors.New("testString error 10")
	}
	n64, err = c.Incrby(nKey, 6)
	if err != nil {
		return err
	}
	if n64 != 10 {
		return errors.New("testString error 11")
	}

	if err := c.Set(key, value); err != nil {
		return err
	}

	exist, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("testString error 12")
	}

	delSize, err := c.Del(key, nKey)
	if err != nil {
		return err
	}
	if delSize != 2 {
		return errors.New("testString error 13")
	}

	exist, err = c.Exists(key)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("testString error 14")
	}

	if _, err := c.DO("FLUSHDB"); err != nil {
		return err
	}

	return nil
}

func testHash(c zerocache.Cache) error {
	var (
		key                    = "hello"
		field1, field2, field3 = "field1", "field2", "field3"
		value1, value2, value3 = "value1", "value2", "value3"
	)

	if err := c.HMSet(key, field1, value1, field2, value2); err != nil {
		return err
	}

	if err := c.HSet(key, field3, value3); err != nil {
		return err
	}

	n, err := c.HLen(key)
	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("testHash error 1")
	}

	n, err = c.HDel(key, field1, field2, "field4")
	if err != nil {
		return err
	}
	if n != 2 {
		return errors.New("testHash error 2")
	}

	vs, err := c.HGetAll(key)
	if err != nil {
		return err
	}
	// field-value
	if len(vs) != 2 {
		return errors.New("testHash error 3")
	}

	v, err := c.HGet(key, field3)
	if err != nil {
		return err
	}
	if v != value3 {
		return errors.New("testHash error 4")
	}

	if _, err := c.DO("FLUSHDB"); err != nil {
		return err
	}

	return nil
}

func testList(c zerocache.Cache) error {
	var (
		key            = "key"
		value1, value2 = "value1", "value2"
	)

	n, _ := c.LPush(key, value1)
	if n != 1 {
		return errors.New("testList error 1")
	}
	n, _ = c.RPush(key, value2)
	if n != 2 {
		return errors.New("testList error 2")
	}

	n, _ = c.LLen(key)
	if n != 2 {
		return errors.New("testList error 3")
	}

	v, _ := c.LIndex(key, 0)
	if v != value1 {
		return errors.New("testList error 4")
	}

	vs, _ := c.LRange(key, 0, 100)
	if len(vs) != 2 || vs[0] != value1 || vs[1] != value2 {
		return errors.New("testList error 5")
	}

	if _, err := c.DO("FLUSHDB"); err != nil {
		return err
	}
	return nil
}

func testSet(c zerocache.Cache) error {
	var (
		key        = "key"
		m1, m2, m3 = "m1", "m2", "m3"
	)

	n, _ := c.SAdd(key, m1, m2, m3)
	n2, _ := c.SCard(key)
	if n != n2 {
		return errors.New("testSet error 1")
	}

	r, _ := c.SIsMember(key, m3)
	if !r {
		return errors.New("testSet error 2")
	}
	r, _ = c.SIsMember(key, "m4")
	if r {
		return errors.New("testSet error 3")
	}

	if _, err := c.DO("FLUSHDB"); err != nil {
		return err
	}

	return nil
}

func testSortedSet(c zerocache.Cache) error {
	var (
		key = "key"
	)

	n, _ := c.ZAdd(key, "99", "m1", "60", "m2", "88", "m3")
	n2, _ := c.ZCard(key)
	if n != n2 {
		return errors.New("testSortedSet error 1")
	}

	n3, _ := c.ZCount(key, 80, 100)
	if n3 != 2 {
		return errors.New("testSortedSet error 2")
	}

	// c.DO("FLUSHDB")
	return nil
}
