package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

var cc cache.Cache

func InitCache() {
	cacheConfig := beego.AppConfig.String("cache")
	cc = nil
	if "redis" == cacheConfig {
		initRedis()
	} else if "memory" == cacheConfig {
		initMemory()
	} else {
		initFile()
	}
	if cc != nil {
		beego.Info("Init cache success!!")
	}
}

func initFile() {
	var err error
	cc, err = cache.NewCache("file", `{"CachePath":"./tmp/cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"180"}`)
	if err != nil {
		beego.Error("New file cache error", err)
	}
}

func initMemory() {
	var err error
	cc, err = cache.NewCache("memory", `{"interval":"180"}`)
	if err != nil {
		beego.Error("New memory cache error", err)
	}
}

func initRedis() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()
	key := beego.AppConfig.String("cacherediskey")
	conn := beego.AppConfig.String("cacheredishost")
	password := beego.AppConfig.String("cacheredispass")
	cc, err = cache.NewCache("redis", `{"key":"`+key+`","conn":"`+conn+`","password":"`+password+`"}`)

	if err != nil {
		beego.Error("New redis cache error", err)
	}
}

func SetCache(key string, value interface{}, timeoutSecond int) error {
	data, err := Encode(value)
	if err != nil {
		beego.Error("Set cache error:", err)
		return err
	}
	if cc == nil {
		beego.Error("Set cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			beego.Error("recover cache error:", r)
			//cc = nil
		}
	}()
	timeouts := time.Duration(timeoutSecond) * time.Second
	err = cc.Put(key, data, timeouts)
	if err != nil {
		beego.Error("Set cache error:", err)
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		beego.Error("Get cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			beego.Error("recover cache error:", r)
			//cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		beego.Warn("Get cache warn Cache不存在")
		return nil
	}
	err := Decode(data.([]byte), to)
	if err != nil {
		beego.Error(err)
	}

	return err
}

func DelCache(key string) error {
	if cc == nil {
		beego.Error("Delete cache error cache is nil")
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			beego.Error("recover cache error:", r)
			//cc = nil
		}
	}()

	err := cc.Delete(key)
	if err != nil {
		beego.Error("Delete cache error", err)
		return err
	} else {
		return nil
	}
}

// --------------------
// Encode
// 用gob进行数据编码
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	if len(data) == 0 {
		beego.Info("Decode data is empty")
		return nil
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
