package db

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response1 struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	NewVersion string `json:"new_version,omitempty"` // 如果没有传新版本号，可以不返回该字段
}

// UpdateData 定义版本更新数据
type UpdateData struct {
	CurrentVersion string `json:"current_version"`
	NewestVersion  string `json:"newest_version"`
	NeedsUpdate    bool   `json:"needs_update"`
}

// Response 用于统一响应结构体
type Response struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Data    UpdateData `json:"data"`
}

// NeedsUpdate 检查更新
// @Summary 检查更新
// @Description 根据客户端提供的当前版本号，检查是否需要更新。
// @Tags update
// @Accept  json
// @Produce  json
// @Param current_version query string false "当前版本号" default(v1.0.0)
// @Success 200 {object} Response "成功返回最新版本信息"
// @Failure 400 {object} Response "请求参数错误"
// @Router /db/echoDB/check-update [get]
func NeedsUpdate(context *gin.Context) {
	// 获取客户端传递的当前版本号，默认是 v1.0.0
	currentVersion := context.DefaultQuery("current_version", "v1.0.0")

	// 获取服务器上的最新版本号
	newestVersion := GetNewestVersion()

	// 比较版本
	needsUpdate := currentVersion != newestVersion

	// 构建响应数据
	response := Response{
		Code:    "200",
		Message: "success",
		Data: UpdateData{
			CurrentVersion: currentVersion,
			NewestVersion:  newestVersion,
			NeedsUpdate:    needsUpdate,
		},
	}

	// 返回结果
	context.JSON(http.StatusOK, response)
}

// UpdateVersion 更新应用版本
// @Summary 更新应用版本
// @Description 更新服务器上记录的最新版本号。
// @Tags update
// @Accept  json
// @Produce  json
// @Param new_version body string true "新的版本号"
// @Success 200 {object} Response "版本更新成功"
// @Failure 400 {object} Response "无效的输入数据"
// @Router /db/echoDB/update-version [post]
func UpdateVersion(context *gin.Context) {
	var json struct {
		NewVersion string `json:"new_version" binding:"required"`
	}

	// 绑定请求体
	if err := context.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回400错误
		context.JSON(http.StatusBadRequest, Response{
			Code:    "400",
			Message: "Invalid input",
		})
		return
	}

	// 更新最新版本
	SetNewestVersion(json.NewVersion)

	// 返回成功消息
	context.JSON(http.StatusOK, Response1{
		Code:       "200",
		Message:    "Version updated successfully",
		NewVersion: json.NewVersion,
	})
}
