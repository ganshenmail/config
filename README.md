# Go 配置管理包

一个线程安全的键值对配置管理库，用于Go应用程序。

## 特性

- **线程安全操作**：所有方法都通过RWMutex保护
- **灵活访问**：获取值时支持默认值回退
- **批量操作**：一次性获取所有配置
- **简单API**：易于集成到任何Go项目中

## 安装

```bash
go get github.com/ganshenmail/config
```

## 快速开始

```go
package main

import (
	"fmt"
	"log"
	"github.com/ganshenmail/config"
)

func main() {
	// 创建新的配置实例
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 从文件加载配置
	err = cfg.LoadFromFile("config.ini")
	if err != nil {
		log.Printf("警告: %v, 使用默认值", err)
	}

	// 获取值
	port := cfg.GetWithDefault("server.port", "8080")
	env := cfg.GetWithDefault("environment", "development")

	fmt.Printf("正在%s模式下启动服务，端口%s\n", port, env)
}
```

## API参考


主要方法:

| 方法 | 描述 |
|--------|-------------|
| `NewConfig()` | 创建新的Config实例 |
| `LoadFromFile(filename)` | 从文件加载配置 |
| `Get(key)` | 根据键获取值 |
| `GetWithDefault(key, defaultValue)` | 获取值，支持默认值回退 |
| `Set(key, value)` | 设置键值对 |
| `SaveToFile(filename)` | 保存配置到文件 |

## 文件格式

配置文件应使用简单的键值对格式:

```ini
# 示例 config.ini
server.port = 8080
environment = production
db.host = localhost
```
