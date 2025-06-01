// Package config 提供线程安全的键值对配置管理系统
//
// 特性:
// - 从文件加载配置(key=value格式)
// - 保存配置到文件
// - 线程安全的Get/Set操作
// - 支持默认值
// - 批量操作(GetAll)
//
// 示例:
//   cfg, err := config.NewConfig()
//   if err != nil {
//       log.Fatal(err)
//   }
//   err = cfg.LoadFromFile("config.ini")
//   port := cfg.GetWithDefault("server.port", "8080")
package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
)

// Config 表示线程安全的键值对配置存储
// 提供加载、保存和操作配置值的方法
// 所有操作都通过RWMutex保护以实现并发访问
type Config struct {
	data  map[string]string
	mutex sync.RWMutex // 保证并发安全
}

// NewConfig 创建并返回新的Config实例
// 返回:
// - *Config: 指向新Config实例的指针
// - error: 初始化错误(如果有)
func NewConfig() (*Config, error) {
	return &Config{
		data: make(map[string]string),
	}, nil
}

// LoadFromFile 从key=value格式的文件加载配置
// 跳过空行和以#开头的行(注释)
// 参数:
// - filename: 配置文件路径
// 返回:
// - error: 文件操作或解析错误(如果有)
func (c *Config) LoadFromFile(filename string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data == nil {
		c.data = make(map[string]string)
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue // 跳过空行和注释
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			c.data[key] = value
		}
	}

	return scanner.Err()
}

// Get 根据键获取配置值
// 参数:
// - key: 要查找的配置键
// 返回:
// - string: 键存在时返回对应值，否则返回空字符串
func (c *Config) Get(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data[key]
}

// GetWithDefault 获取配置值，支持默认值回退
// 参数:
// - key: 要查找的配置键
// - defaultValue: 键不存在时返回的默认值
// 返回:
// - string: 键存在时返回对应值，否则返回defaultValue
func (c *Config) GetWithDefault(key, defaultValue string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.data[key]; ok {
		return val
	}
	return defaultValue
}

// Set 存储配置值
// 参数:
// - key: 配置键
// - value: 要存储的值
// 返回:
// - error: 当key为空时返回错误
func (c *Config) Set(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return nil
}

// Has 检查配置键是否存在
// 参数:
// - key: 要检查的配置键
// 返回:
// - bool: 键存在时返回true
func (c *Config) Has(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data[key]
	return ok
}

// Delete 删除配置键值对
// 参数:
// - key: 要删除的配置键
func (c *Config) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// GetAll 返回所有配置键值对的副本
// 返回:
// - map[string]string: 所有配置数据的副本
func (c *Config) GetAll() map[string]string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	copy := make(map[string]string, len(c.data))
	for k, v := range c.data {
		copy[k] = v
	}
	return copy
}

// SaveToFile 将所有配置以key=value格式保存到文件
// 参数:
// - filename: 目标文件路径
// 返回:
// - error: 文件操作错误(如果有)
func (c *Config) SaveToFile(filename string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range c.data {
		_, err := writer.WriteString(key + " = " + value + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}