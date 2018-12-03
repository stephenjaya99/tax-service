package handler

import (
	c "gitlab.com/stephenjaya99/tax-service/controller"
	m "gitlab.com/stephenjaya99/tax-service/model"

	"github.com/gin-gonic/gin"
)

// handler holds the structure for Handler
type handler struct {
	controller c.Controller
}

// Handler holds the contract for Handler
type Handler interface {
	// Ping should handle healthcheck in top level routing
	Ping(*gin.Context)
	RetrieveAllTaxes(*gin.Context)
	CreateTax(*gin.Context)
}

// response define default structure of handler response
type response struct {
	// Meta contains metadata of the response
	Meta meta
	// Body contains content of response
	Body interface{}
}

// New is a function for creating handler
func New(controller c.Controller) Handler {
	return &handler{
		controller: controller,
	}
}

// Ping if a function for handling healthcheck in top level routing
func (h *handler) Ping(ctx *gin.Context) {
	select {
	case <-ctx.Request.Context().Done():
		m := metaContextTimeout
		resp := &response{m, nil}
		ctx.JSON(408, resp)
		return
	default:
	}

	resp := &response{metaSuccess, "Tax Service is running!"}
	ctx.JSON(200, resp)

	return
}

// CreateTax is a function for handling create tax request in top level routing
func (h *handler) CreateTax(ctx *gin.Context) {
	select {
	case <-ctx.Request.Context().Done():
		m := metaContextTimeout
		resp := &response{m, nil}
		ctx.JSON(408, resp)
		return
	default:
	}

	var taxRequest c.TaxRequest
	err := ctx.ShouldBindJSON(&taxRequest)
	if err != nil {
		m := metaJSONDecodeError
		m.Error = err.Error()
		resp := &response{m, nil}
		ctx.JSON(400, resp)
		return
	}

	var tax m.Tax
	tax, err = h.controller.CreateTax(ctx, taxRequest)
	if err != nil {
		m := metaControllerError
		m.Error = err.Error()
		resp := &response{m, nil}
		ctx.JSON(500, resp)
		return
	}

	resp := &response{metaSuccess, tax}
	ctx.JSON(201, resp)

	return
}

// RetrieveAllTaxes is a function for getting all tax request in top level routing
func (h *handler) RetrieveAllTaxes(ctx *gin.Context) {
	select {
	case <-ctx.Request.Context().Done():
		m := metaContextTimeout
		resp := &response{m, nil}
		ctx.JSON(408, resp)
		return
	default:
	}

	taxes, err := h.controller.RetrieveAllTaxes(ctx)
	if err != nil {
		m := metaControllerError
		m.Error = err.Error()
		resp := &response{m, nil}
		ctx.JSON(500, resp)
		return
	}

	resp := &response{metaSuccess, taxes}
	ctx.JSON(200, resp)

	return
}
