// Package config provides a thread-safe key-value configuration management system.
//
// Features:
// - Load configuration from file (key=value format)
// - Save configuration to file
// - Get/Set values with thread safety
// - Default values support
// - Bulk operations (GetAll)
//
// Example:
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

// Config represents a thread-safe key-value configuration store.
// It provides methods to load, save, and manipulate configuration values.
// All operations are protected by a RWMutex for concurrent access.
type Config struct {
	data  map[string]string
	mutex sync.RWMutex // 保证并发安全
}

// NewConfig creates and returns a new Config instance.
// Returns:
// - *Config: pointer to new Config instance
// - error: any initialization error
func NewConfig() (*Config, error) {
	return &Config{
		data: make(map[string]string),
	}, nil
}

// LoadFromFile loads configuration from a file in key=value format.
// Skips empty lines and lines starting with # (comments).
// Args:
// - filename: path to configuration file
// Returns:
// - error: any file operation or parsing error
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

// Get retrieves a configuration value by key.
// Args:
// - key: configuration key to lookup
// Returns:
// - string: value if key exists, empty string otherwise
func (c *Config) Get(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data[key]
}

// GetWithDefault retrieves a configuration value with fallback to default.
// Args:
// - key: configuration key to lookup
// - defaultValue: value to return if key doesn't exist
// Returns:
// - string: value if key exists, defaultValue otherwise
func (c *Config) GetWithDefault(key, defaultValue string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.data[key]; ok {
		return val
	}
	return defaultValue
}

// Set stores a configuration value.
// Args:
// - key: configuration key
// - value: value to store
// Returns:
// - error: if key is empty
func (c *Config) Set(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return nil
}

// Has checks if a configuration key exists.
// Args:
// - key: configuration key to check
// Returns:
// - bool: true if key exists
func (c *Config) Has(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data[key]
	return ok
}

// Delete removes a configuration key-value pair.
// Args:
// - key: configuration key to remove
func (c *Config) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// GetAll returns a copy of all configuration key-value pairs.
// Returns:
// - map[string]string: copy of all configuration data
func (c *Config) GetAll() map[string]string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	copy := make(map[string]string, len(c.data))
	for k, v := range c.data {
		copy[k] = v
	}
	return copy
}

// SaveToFile saves all configuration to a file in key=value format.
// Args:
// - filename: path to destination file
// Returns:
// - error: any file operation error
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