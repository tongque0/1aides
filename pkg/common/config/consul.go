package config

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/consul/api"
)

// ConsulConfig 是一个封装了 Consul 客户端的配置管理器。
type ConsulConfig struct {
	client *api.Client
}

// NewConsulConfig 创建并返回一个新的 ConsulConfig 实例。
func NewConsulConfig(address string) (*ConsulConfig, error) {
	// 配置 Consul 客户端。
	config := api.DefaultConfig()
	config.Address = address

	// 创建 Consul 客户端。
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}

	return &ConsulConfig{
		client: client,
	}, nil
}

// LoadConfig 从 Consul 中读取配置并解析到指定的结构体中。
func (c *ConsulConfig) LoadConfig(key string, config interface{}) error {
	kv := c.client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return fmt.Errorf("从 Consul 获取配置失败: %v", err)
	}
	if pair == nil {
		return fmt.Errorf("未找到键: %s", key)
	}

	err = yaml.Unmarshal(pair.Value, config)
	if err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	return nil
}

// WatchConfig 监控 Consul 中指定键的配置变化，并在配置变化时调用回调函数。
func (c *ConsulConfig) WatchConfig(key string, config interface{}, onChange func()) error {
	var lastIndex uint64

	for {
		kv := c.client.KV()
		pair, meta, err := kv.Get(key, &api.QueryOptions{
			WaitIndex: lastIndex,
			WaitTime:  10 * time.Second,
		})
		if err != nil {
			log.Printf("监控 Consul 配置失败: %v", err)
			time.Sleep(time.Second * 5) // 遇到错误时，等待一段时间再重试
			continue
		}

		// 检查配置是否有更新。
		if meta.LastIndex > lastIndex {
			lastIndex = meta.LastIndex
			err = yaml.Unmarshal(pair.Value, config)
			if err != nil {
				log.Printf("解析配置失败: %v", err)
			} else {
				log.Printf("配置发生变化: %+v", config)
				onChange() // 配置发生变化时调用回调函数。
			}
		}
	}
}
