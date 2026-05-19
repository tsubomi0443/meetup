package interfaces

import "github.com/labstack/echo/v5"

// setStaticHandler は静的ファイル・Favicon.icoのルートを登録する。
func (r *Router) setStaticHandler() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, r.e.File("/favicon.ico", "static/images/favicon.ico"))
	routeInfos = append(routeInfos, r.e.Static("/static", "static"))
	return
}
