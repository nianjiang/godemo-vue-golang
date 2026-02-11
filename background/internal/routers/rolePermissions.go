package routers

import (
	"github.com/gin-gonic/gin"

	"godemo/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		rolePermissionsRouter(group, handler.NewRolePermissionsHandler())
	})
}

func rolePermissionsRouter(group *gin.RouterGroup, h handler.RolePermissionsHandler) {
	g := group.Group("/rolePermissions")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)          // [post] /api/v1/rolePermissions
	g.DELETE("/:roleID", h.DeleteByRoleID) // [delete] /api/v1/rolePermissions/:roleID
	g.PUT("/:roleID", h.UpdateByRoleID)    // [put] /api/v1/rolePermissions/:roleID
	g.GET("/:roleID", h.GetByRoleID)       // [get] /api/v1/rolePermissions/:roleID
	g.POST("/list", h.List)        // [post] /api/v1/rolePermissions/list
}
