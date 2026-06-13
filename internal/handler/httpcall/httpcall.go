// Package httpcall HTTP 调用示例。演示如何用 net/http 请求外部 API。
package httpcall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/logger"
	"github.com/wxsimon2022/simonStu/internal/response"
)

// HttpCallRequest 外部 HTTP 调用请求参数。
type HttpCallRequest struct {
	URL    string `json:"url" binding:"required"` // 目标 URL
	Method string `json:"method"`                 // GET / POST，默认 GET
	Body   any    `json:"body"`                   // POST 时发送的 JSON body
	Param  string `json:"param"`                  // URL 查询参数值（如 ?q=xxx）
}

// HttpCallResult 外部 HTTP 调用返回结果。
type HttpCallResult struct {
	URL        string      `json:"url"`
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       any         `json:"body"`
	DurationMs int64       `json:"duration_ms"`
}

// HttpCall 发送 HTTP 请求调用外部 API，返回响应状态码、响应头和响应体。
//
// 请求示例：
//
//	GET  /commonApi/http/call  ?url=https://httpbin.org/get&param=hello
//	POST /commonApi/http/call  {"url":"https://httpbin.org/post","method":"POST","body":{"name":"simon"}}
func HttpCall(c *gin.Context) {
	start := time.Now()

	// 解析参数：支持 GET query 和 POST JSON
	var req HttpCallRequest

	if c.Request.Method == http.MethodGet {
		req.URL = c.DefaultQuery("url", "")
		req.Method = c.DefaultQuery("method", "GET")
		req.Param = c.DefaultQuery("param", "")
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Errorf(c, "HttpCall 参数解析失败: %v", err)
			response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
			return
		}
	}

	if req.URL == "" {
		logger.Errorf(c, "HttpCall url 为空")
		response.Error(c, http.StatusBadRequest, "参数 url 不能为空")
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	// 构造请求 body
	var bodyReader io.Reader
	if req.Body != nil {
		b, err := json.Marshal(req.Body)
		if err != nil {
			logger.Errorf(c, "HttpCall body 序列化失败: %v", err)
			response.Error(c, http.StatusBadRequest, "body 序列化失败")
			return
		}
		bodyReader = bytes.NewReader(b)
	}

	// 创建 HTTP 请求（继承 Gin 的 context，请求取消时自动取消）
	httpReq, err := http.NewRequestWithContext(c.Request.Context(), req.Method, req.URL, bodyReader)
	if err != nil {
		logger.Errorf(c, "HttpCall 创建请求失败 url=%s err=%v", req.URL, err)
		response.Error(c, http.StatusBadRequest, "创建请求失败: "+err.Error())
		return
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// 发送请求（带 10 秒超时）
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Errorf(c, "HttpCall 请求失败 url=%s err=%v", req.URL, err)
		response.Error(c, http.StatusBadGateway, "请求外部接口失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应 body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(c, "HttpCall 读取响应失败 url=%s err=%v", req.URL, err)
		response.Error(c, http.StatusBadGateway, "读取响应失败: "+err.Error())
		return
	}

	// 尝试将响应体解析为 JSON，否则返回原始字符串
	var parsedBody any
	if err := json.Unmarshal(respBody, &parsedBody); err != nil {
		parsedBody = string(respBody)
	}

	duration := time.Since(start).Milliseconds()

	logger.Infof(c, "HttpCall %s %s → %d (%dms)", req.Method, req.URL, resp.StatusCode, duration)

	response.Success(c, HttpCallResult{
		URL:        req.URL,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       parsedBody,
		DurationMs: duration,
	})
}

// HttpCallLocal 演示调用本项目的另一个接口（内部调用，不需经过网络）。
func HttpCallLocal(c *gin.Context) {
	// 用本项目的 base URL（适合 docker 内或本地调试时调用其他接口）
	baseURL := fmt.Sprintf("http://localhost:%s", c.Request.Host)
	targetURL := baseURL + "/commonApi/hello"

	start := time.Now()

	ctx := c.Request.Context()
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Errorf(c, "HttpCallLocal 调用失败 url=%s err=%v", targetURL, err)
		response.Error(c, http.StatusBadGateway, "内部调用失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	duration := time.Since(start).Milliseconds()

	logger.Infof(c, "HttpCallLocal %s → %d (%dms)", targetURL, resp.StatusCode, duration)

	response.Success(c, gin.H{
		"target_url":  targetURL,
		"status_code": resp.StatusCode,
		"body":        string(body),
		"duration_ms": duration,
	})
}
