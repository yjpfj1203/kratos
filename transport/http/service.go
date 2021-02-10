package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

type methodHandler func(srv interface{}, ctx context.Context, req *http.Request, m middleware.Middleware) (out interface{}, err error)

// MethodDesc represents a Proto service's method specification.
type MethodDesc struct {
	Path    string
	Method  string
	Handler methodHandler
}

// ServiceDesc represents a Proto service's specification.
type ServiceDesc struct {
	ServiceName string
	Methods     []MethodDesc
	Metadata    interface{}
}

// ServiceRegistrar wraps a single method that supports service registration.
type ServiceRegistrar interface {
	RegisterService(desc *ServiceDesc, impl interface{})
}

// RegisterService .
func (s *Server) RegisterService(desc *ServiceDesc, impl interface{}) {
	for _, m := range desc.Methods {
		h := m.Handler
		s.router.HandleFunc(m.Path, func(res http.ResponseWriter, req *http.Request) {
			out, err := h(impl, req.Context(), req, s.opts.middleware)
			if err != nil {
				s.opts.errorEncoder(res, req, err)
				return
			}
			if err := s.opts.responseEncoder(res, req, out); err != nil {
				s.opts.errorEncoder(res, req, err)
			}
		}).Methods(m.Method)
	}
}
