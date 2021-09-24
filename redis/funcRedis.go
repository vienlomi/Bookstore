package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"project_v3/data"
	"time"

	"github.com/go-redis/redis"
)

//func set key - value into redis with time expire
func ClientSet(key string, obj interface{}, time time.Duration) error {
	value, err := json.Marshal(&obj)
	if err != nil {
		fmt.Println("can not marshal obj")
		return err
	}
	fmt.Println(value)
	err = Client.Set(key, value, time).Err()
	if err != nil {
		fmt.Println("can not set obj into redis", err)
		return err
	}
	return nil
}

//func get value of key in redis
func ClientGet(key string, obj interface{}) error {
	value, err := Client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
		return err
	}
	if err != nil {
		fmt.Println("can not get value ", err)
		return err
	}
	err = json.Unmarshal([]byte(value), &obj)
	if err != nil {
		fmt.Println("can not unmarshal ", err)
		return err
	}
	return nil
}

// func RPush a new list with key - value is slice of struct and expire time forever
func ClientRPushProduct(key string, obj []data.ProductLite) error {
	ClientDelete(key)
	fmt.Println(key)
	var values []interface{}
	for _, val := range obj {
		js, err := json.Marshal(&val)
		if err != nil {
			fmt.Println("can not marshal obj [] ", err)
			return err
		}
		values = append(values, js)
	}

	for _, val := range values {
		rs, err := Client.RPush(key, val).Result()
		if err != nil || rs == 0 {
			fmt.Println("can not LPush redis ", err)
		}
	}
	//rs, err := Client.Expire(CtxRedis, key, time).Result()
	//if err != nil {
	//	fmt.Println("set expire error: ", err)
	//}
	//fmt.Println("set expire: ", rs)
	//timecl := Client.TTL(CtxRedis, key).Val()
	//fmt.Println("time TLL", timecl)
	return nil
}

// func RPush a new list with key - value is slice of struct and expire time be set
func ClientRPushProductExpire(key string, obj []data.ProductLite, time time.Duration) error {
	ClientDelete(key)
	fmt.Println(key)
	var values []interface{}
	for _, val := range obj {
		js, err := json.Marshal(&val)
		if err != nil {
			fmt.Println("can not marshal obj [] ", err)
			return err
		}
		values = append(values, js)
	}

	for _, val := range values {
		rs, err := Client.RPush(key, val).Result()
		if err != nil || rs == 0 {
			fmt.Println("can not LPush redis ", err)
		}
	}
	rs, err := Client.Expire(key, time).Result()
	if err != nil {
		fmt.Println("set expire error: ", err)
	}
	fmt.Println("set expire: ", rs)
	timecl := Client.TTL(key).Val()
	fmt.Println("time TLL", timecl)
	return nil
}

// func add slice's obj to key (available or not available) with time expire
func ClientAddRPushProduct(key string, obj []data.ProductLite, time time.Duration) error {
	var values []interface{}
	for _, val := range obj {
		js, err := json.Marshal(&val)
		if err != nil {
			fmt.Println("can not marshal obj [] ", err)
			return err
		}
		values = append(values, js)
	}

	Client.Expire(key, time)

	for _, val := range values {
		err := Client.RPush(key, val)
		if err != nil {
			fmt.Println("can not LPush redis ", err)
		}
	}
	return nil
}

// func LRange list in redis, return slice of obj
func ClientLRangProduct(key string, obj *[]data.ProductLite) error {
	fmt.Println(key)

	var tmp data.ProductLite
	value, err := Client.LRange(key, 0, -1).Result()
	fmt.Println(value)
	if err != nil || len(value) == 0 {
		fmt.Println("can not get value ", err)
		return errors.New("can not get value list redis")
	}
	for _, val := range value {
		err = json.Unmarshal([]byte(val), &tmp)
		if err != nil {
			fmt.Println("can not unmarshal obj ", err)
			return err
		}
		*obj = append(*obj, tmp)
	}
	return nil
}

//func delete key in redis
func ClientDelete(key string) {
	count := Client.Del(key)
	if count.Val() > 0 {
		fmt.Println("deleted key in redis")
	} else {
		fmt.Println("no key deleted in redis")
	}
}

func ClientDeleteAllKey() {
	Client.FlushAll()
}
