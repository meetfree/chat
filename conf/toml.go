package conf

import (
	"github.com/BurntSushi/toml"
)

// Config 对应配置文件结构
type Toml struct {
	Listen  string                 `toml:"listen"`
	Dict    string                 `toml:"dict"`
	NatAddr string				   `toml:"nat"`
	Redis   map[string]RedisServer `toml:"redis"`
	Mongo   map[string]MongoServer `toml:"mongo"`
}

// UnmarshalConfig 解析toml配置
func UnmarshalConfig(file string) (*Toml, error) {
	c := &Toml{}
	if _, err := toml.DecodeFile(file, c); err != nil {
		return c, err
	}
	return c, nil
}

// RedisServerConf 获取数据库配置
func (t *Toml) RedisConfig(key string) RedisServer {
	return t.Redis[key]
}

// 获取数据库配置
func (t *Toml) MongoConfig(key string) MongoServer {
	return t.Mongo[key]
}

// 监听地址
func (t *Toml) GetListenAddr() string {
	return t.Listen
}

// NatS地址
func (t *Toml) NatConfig() string {
	return t.NatAddr
}

// RedisServer 表示 redis 服务器配置
type RedisServer struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// RedisServer 表示 redis 服务器配置
type MongoServer struct {
	Addr     string `toml:"addr"`
	Database string `toml:"database"`
}
