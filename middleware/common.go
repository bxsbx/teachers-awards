package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"teachers-awards/common/tracer"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"time"
)

const (
	REQ_ID = "req_id"
)

// 做一些公共处理，比如上下文，记录一些日志信息。
func Common() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := context.Background()
		reqId := util.GetRandUUID()
		ctx = context.WithValue(ctx, REQ_ID, reqId)
		//开启链路
		ctx = tracer.StarTracerSpan(ctx, c.Request)
		tracer.SetTagSpan(ctx, REQ_ID, reqId)
		req := c.Request
		body, _ := io.ReadAll(req.Body)
		if len(body) > 0 {
			//body的数据被读取出来后需要重新设置回去
			req.Body = io.NopCloser(bytes.NewReader(body))
		}
		global.SetContext(c, ctx)
		c.Next()
		latency := time.Since(start) //记录耗时

		fields := []zap.Field{
			zap.String("ip", c.ClientIP()),
			zap.String("scheme", req.Proto),
			zap.String("remoteAddr", req.RemoteAddr),
			zap.String("url", req.URL.String()),
			zap.String("method", req.Method),
			zap.Any("header", req.Header),
			zap.Any("form", req.Form),
			zap.String("body", string(body)),
			zap.String("latency", fmt.Sprintf("%dms", latency.Milliseconds())),
		}
		value, exists := c.Get("errStack")
		if exists {
			fields = append(fields, zap.Any("errStack", value))
			global.Logger.Error("Request Error", fields...)
		} else {
			global.Logger.Info("Request Info", fields...)
		}
		//释放链路
		tracer.FinishSpan(global.GetContext(c))
	}
}
