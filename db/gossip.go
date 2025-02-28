package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GossipEngine struct {
	peers       []string
	port        int
	nodeID      string
	lastGossip  time.Time
	gossipState string
}

// GossipState 定义节点的状态结构
type GossipState struct {
	NodeID     string `json:"node_id"`
	LastGossip int64  `json:"last_gossip"`
	State      string `json:"state"`
}

// NewGossipEngine 创建一个新的Gossip引擎实例
func NewGossipEngine(peers []string, port int, nodeID string) *GossipEngine {
	return &GossipEngine{
		peers:       peers,
		port:        port,
		nodeID:      nodeID,
		lastGossip:  time.Now(),
		gossipState: fmt.Sprintf("Node %s is running", nodeID),
	}
}

// 启动Gossip服务
func (g *GossipEngine) StartGossipServer() {
	router := gin.Default()

	// Gossip 接口，用于接收来自其他节点的 Gossip 数据
	router.POST("/gossip", func(c *gin.Context) {
		var receivedState GossipState
		if err := c.ShouldBindJSON(&receivedState); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gossip data"})
			return
		}

		// 更新节点的 Gossip 状态
		g.lastGossip = time.Unix(receivedState.LastGossip, 0)
		g.gossipState = receivedState.State

		fmt.Printf("Received gossip from node %s: %s\n", receivedState.NodeID, receivedState.State)

		c.JSON(http.StatusOK, gin.H{
			"message": "Gossip received successfully",
		})
	})

	// 启动 HTTP 服务，使用新的端口
	go func() {
		err := router.Run(fmt.Sprintf(":%d", g.port)) // 启动服务器，确保 g.port 是修改后的端口
		if err != nil {
			fmt.Printf("Failed to start gossip server on port %d: %v\n", g.port, err)
		}
	}()
}

// Gossip 向其他节点传播数据
func (g *GossipEngine) Gossip() {
	for _, peer := range g.peers {
		fmt.Printf("Gossiping to peer: %s\n", peer)

		// 向其他节点传播 Gossip 数据
		gossipData := GossipState{
			NodeID:     g.nodeID,
			LastGossip: time.Now().Unix(),
			State:      g.gossipState,
		}

		// 向其他节点发送 Gossip 数据
		err := sendGossip(peer, gossipData)
		if err != nil {
			fmt.Printf("Failed to gossip to node %s: %v\n", peer, err)
		}
	}
}

// 向指定节点发送 Gossip 数据
func sendGossip(peer string, data GossipState) error {
	// 构建 HTTP 请求
	url := fmt.Sprintf("http://%s/gossip", peer)
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置请求超时时间
	}

	// 将 GossipState 编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal gossip data: %w", err)
	}

	// 使用 bytes.Buffer 创建请求体
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send gossip to %s: %w", peer, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from %s: %v", peer, resp.StatusCode)
	}

	return nil
}
