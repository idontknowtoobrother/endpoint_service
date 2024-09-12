package endpoint

import (
	"errors"
	"fmt"

	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/idontknowtoobrother/practice_go_hexagonal/constraints"
	"github.com/idontknowtoobrother/practice_go_hexagonal/services/endpoint"
)

type endpointHandler struct {
	srv endpoint.EndpointService
}

func NewEndpointHandler(srv endpoint.EndpointService) *endpointHandler {
	return &endpointHandler{
		srv: srv,
	}
}

func (h *endpointHandler) GetEndpoints(c *gin.Context) {
	endpoints, err := h.srv.GetEndpoints()
	if err != nil {
		utils.ErrorResponse(c, 500, err)
		return
	}
	utils.SuccessResponse(c, 200, endpoints)
}

func (h *endpointHandler) GetEndpoint(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		utils.ErrorResponse(c, 400, errors.New(fmt.Sprintf(constraints.ErrInvalidParameter, "uuid")))
		return
	}
	endpoint, err := h.srv.GetEndpoint(uuid)
	if err != nil {
		utils.ErrorResponse(c, 500, err)
		return
	}
	utils.SuccessResponse(c, 200, endpoint)
}
