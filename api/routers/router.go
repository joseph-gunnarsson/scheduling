package routers

import (
	"net/http"

	"github.com/joseph-gunnarsson/scheduling/api/handlers"
	"github.com/joseph-gunnarsson/scheduling/api/middleware"
)

func Routers(handler *handlers.BaseHandler, mm *middleware.MiddlewareManager) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/register/", middleware.MultipleMiddleware(handler.CreateUserHandler, mm.ErrorHandlerMiddleware))
	mux.HandleFunc("POST /user/updatepassword/{id}/", middleware.MultipleMiddleware(handler.UpdatePassword, mm.AuthMiddleware, mm.ErrorHandlerMiddleware))

	mux.HandleFunc("POST /user/login/", middleware.MultipleMiddleware(handler.LoginHandler, mm.ErrorHandlerMiddleware))

	mux.HandleFunc("POST /group/", middleware.MultipleMiddleware(handler.CreateGroupHandler, mm.AuthMiddleware, mm.ErrorHandlerMiddleware))
	mux.HandleFunc("DELETE /group/{id}/", middleware.MultipleMiddleware(handler.DeleteGroupHandler, mm.ErrorHandlerMiddleware, mm.AuthMiddleware, mm.GroupPermissionMiddleware))
	mux.HandleFunc("PATCH /group/{id}/", middleware.MultipleMiddleware(handler.PatchGroupHandler, mm.ErrorHandlerMiddleware, mm.AuthMiddleware, mm.GroupPermissionMiddleware))
	mux.HandleFunc("PUT /group/{id}/", middleware.MultipleMiddleware(handler.UpdateGroupHandler, mm.ErrorHandlerMiddleware, mm.AuthMiddleware, mm.GroupPermissionMiddleware))

	mux.HandleFunc("GET /user/{id}/group/", middleware.MultipleMiddleware(handler.GetGroupsByOwnerHandler, mm.ErrorHandlerMiddleware, mm.AuthMiddleware))

	return mux
}
