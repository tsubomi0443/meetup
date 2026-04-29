package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (hm *HandlerManager) SetSSEHandler() (routeInfos []echo.RouteInfo) {
	// SSEサーバの開始（/sse）
	go hm.hub.Run()
	// SSEの定期送信処理の開始（時刻表示など）
	go hm.hub.RunSSE(hm.db)

	group := hm.e.Group("", GetJWTConfig())
	routeInfos = append(routeInfos, group.GET("/sse", hm.sseHandler()))
	return
}

// sseHandler はSSEストリームを提供するEchoハンドラを返す。
func (hm *HandlerManager) sseHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		client := &Client{send: make(chan Event, 64)}
		hm.hub.Register <- client
		defer func() {
			hm.hub.Unregister <- client
		}()

		resWriter := c.Response()
		resWriter.Header().Set("Content-Type", "text/event-stream")
		resWriter.Header().Set("Cache-Control", "no-cache")
		resWriter.Header().Set("Connection", "keep-alive")
		resWriter.Header().Set("Transfer-Encoding", "chunked")
		resWriter.WriteHeader(http.StatusOK)

		flusher, ok := resWriter.(http.Flusher)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
		}

		ctx := c.Request().Context()
		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-client.send:
				if !ok {
					return nil
				}
				if event.ID != "" {
					if _, err := fmt.Fprintf(resWriter, "id: %s\n", event.ID); err != nil {
						return err
					}
				}
				if event.Event != "" {
					if _, err := fmt.Fprintf(resWriter, "event: %s\n", event.Event); err != nil {
						return err
					}
				}
				if _, err := fmt.Fprintf(resWriter, "data: %s\n\n", event.Data); err != nil {
					return err
				}
				flusher.Flush()
			}
		}
	}
}
