package db

import (
	"fmt"
	"sync"
)

// GlobalVersion 是一个全局版本管理结构
var GlobalVersion struct {
	newestVersion string
	mutex         sync.RWMutex
}

// SetNewestVersion 设置最新版本号
func SetNewestVersion(version string) {
	GlobalVersion.mutex.Lock()
	defer GlobalVersion.mutex.Unlock()
	GlobalVersion.newestVersion = version
	fmt.Printf("Newest version set to: %s\n", version)
}

// GetNewestVersion 获取当前最新版本号
func GetNewestVersion() string {
	GlobalVersion.mutex.RLock()
	defer GlobalVersion.mutex.RUnlock()
	return GlobalVersion.newestVersion
}
