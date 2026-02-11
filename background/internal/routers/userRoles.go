package routers

import (
	"github.com/gin-gonic/gin"

	"godemo/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		userRolesRouter(group, handler.NewUserRolesHandler())
	})
}

func userRolesRouter(group *gin.RouterGroup, h handler.UserRolesHandler) {
	g := group.Group("/userRoles")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)                  // [post] /api/v1/userRoles
	g.DELETE("/:userID", h.DeleteByUserID) // [delete] /api/v1/userRoles/:userID
	g.PUT("/:userID", h.UpdateByUserID)    // [put] /api/v1/userRoles/:userID
	g.GET("/:userID", h.GetByUserID)       // [get] /api/v1/userRoles/:userID
	g.POST("/list", h.List)                // [post] /api/v1/userRoles/list
}
