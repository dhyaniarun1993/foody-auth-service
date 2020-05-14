package controllers

import (
	"github.com/gorilla/mux"
)

// Controller provides interface for controller
type Controller interface {
	LoadRoutes(*mux.Router)
}
