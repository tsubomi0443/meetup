package interfaces

import (
	"fmt"
	"net/http"

	"meetup/internal/interfaces/sse"

	"github.com/labstack/echo/v5"
)

func (r *Router) setSSEHandler() (routeInfos []echo.RouteInfo) {
	go r.deps.Hub.Run()
	go r.deps.Hub.RunSSE()

	group := r.e.Group("", GetJWTConfig())
	routeInfos = append(routeInfos, group.GET("/sse", r.sseHandler()))
	return
}

func (r *Router) sseHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		client := &sse.Client{Send: make(chan sse.Event, 64)}
		r.deps.Hub.Register <- client
		defer func() {
			r.deps.Hub.Unregister <- client
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
			case event, ok := <-client.Send:
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
