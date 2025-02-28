package db

import (
	"echoDB/config"
	"fmt"
	"sync"
	"time"
)

// Item 表示数据库中的一项数据
type Item struct {
	Value        interface{} // 存储的值
	Frequency    int         // 访问频率
	LastAccessed time.Time   // 最后访问时间
	Expiration   time.Time   // 过期时间
}

// EchoDB 是分布式内存数据库的结构
type EchoDB struct {
	data      map[string]*Item // 存储数据
	mutex     sync.RWMutex     // 保护并发访问
	config    *config.Config
	Gossip    *GossipEngine
	bplusTree *BPlusTree    // B+树用于索引
	maxSize   int           // 最大存储条数
	lifetime  time.Duration // 数据过期时间
}

// NewEchoDB 创建一个新的EchoDB实例
func NewEchoDB(config *config.Config) *EchoDB {
	db := &EchoDB{
		data:      make(map[string]*Item),
		config:    config,
		maxSize:   1000,             // 设定最大条目数量
		lifetime:  10 * time.Minute, // 数据过期时间设为10分钟
		bplusTree: NewBPlusTree(3),  // 初始化B+树，假设度为3
	}

	// 根据配置选择一致性算法
	if config.ConsistencyAlgorithm == "Gossip" {
		db.Gossip = NewGossipEngine(config.Gossip.Peers, config.Gossip.Port, config.Gossip.NodeID)
	}

	// 启动定时淘汰任务
	go db.startEvictionProcess()

	return db
}

// Insert 插入或更新数据
func (db *EchoDB) Insert(key string, value interface{}) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// 检查是否需要更新已有的条目
	if item, exists := db.data[key]; exists {
		// 更新数据项的值
		item.Value = value
		// 更新访问频率和最后访问时间
		item.Frequency++
		item.LastAccessed = time.Now()
	} else {
		// 新数据项，设定过期时间
		db.data[key] = &Item{
			Value:        value,
			Frequency:    1,
			LastAccessed: time.Now(),
			Expiration:   time.Now().Add(db.lifetime), // 设置过期时间
		}
	}

	// 更新B+树索引
	db.bplusTree.Insert(key)

	// 如果数据量超出最大存储，进行清理
	if len(db.data) > db.maxSize {
		db.evictData()
	}

	return nil
}

// Delete 删除数据
func (db *EchoDB) Delete(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// 删除数据
	delete(db.data, key)

	// 删除B+树索引
	db.bplusTree.Delete(key)

	return nil
}

// Query 查询数据
func (db *EchoDB) Query(key string) (interface{}, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	item, exists := db.data[key]
	if exists {
		// 更新访问频率和最后访问时间
		item.Frequency++
		item.LastAccessed = time.Now()
	}

	return item, exists
}

// evictData 根据LFU + LRU 策略删除数据
func (db *EchoDB) evictData() {
	var lowestFreqItems []*Item
	var lowestFreq int

	// 找到最低访问频率
	for _, item := range db.data {
		if lowestFreq == 0 || item.Frequency < lowestFreq {
			lowestFreq = item.Frequency
			lowestFreqItems = []*Item{item}
		} else if item.Frequency == lowestFreq {
			lowestFreqItems = append(lowestFreqItems, item)
		}
	}

	// 如果有多个条目访问频率相同，使用LRU策略（最久未访问的条目）
	if len(lowestFreqItems) > 1 {
		// 找到最久未访问的条目
		var oldestItem *Item
		for _, item := range lowestFreqItems {
			if oldestItem == nil || item.LastAccessed.Before(oldestItem.LastAccessed) {
				oldestItem = item
			}
		}
		// 删除最久未访问的条目
		for key, item := range db.data {
			if item == oldestItem {
				delete(db.data, key)
				db.bplusTree.Delete(key)
				break
			}
		}
	} else {
		// 只有一个最低频率的数据，删除它
		for key, item := range db.data {
			if item == lowestFreqItems[0] {
				delete(db.data, key)
				db.bplusTree.Delete(key)
				break
			}
		}
	}
}

// evictExpiredData 清理过期数据
func (db *EchoDB) evictExpiredData() {
	currentTime := time.Now()

	db.mutex.Lock()
	defer db.mutex.Unlock()

	// 遍历所有数据，删除已过期的
	for key, item := range db.data {
		if item.Expiration.Before(currentTime) {
			delete(db.data, key)
			db.bplusTree.Delete(key)
		}
	}
}

// startEvictionProcess 定期检查并淘汰数据
func (db *EchoDB) startEvictionProcess() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer ticker.Stop()

	for range ticker.C {
		// 清理过期数据
		db.evictExpiredData()

		// 如果数据量超出最大存储，进行清理
		if len(db.data) > db.maxSize {
			db.evictData()
		}
	}
}

// RangeQuery 支持范围查询，利用B+树来实现
func (db *EchoDB) RangeQuery(startKey, endKey string) {
	node := db.bplusTree.root
	var result []string

	// 使用 B+树的叶子节点遍历来支持范围查询
	for node != nil {
		for _, key := range node.keys {
			if key >= startKey && key <= endKey {
				result = append(result, key)
			}
		}
		node = node.next
	}

	fmt.Println("Range query result:", result)
}

// PrintIndex 打印B+树的索引结构
func (db *EchoDB) PrintIndex() {
	db.bplusTree.PrintTree(db.bplusTree.root, 0)
}
