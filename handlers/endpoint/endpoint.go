package endpoint

import (
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
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
	uuid := c.Query("uuid")
	if uuid != "" {
		endpoint, err := h.srv.GetEndpoint(uuid)
		if err != nil {
			res := utils.NewErrorResponse(500, err.Error())
			c.JSON(res.StatusCode, res)
			return
		}
		res := utils.NewSuccessResponse(200, "endpoint found.", endpoint)
		c.JSON(res.StatusCode, res)
		return
	}

	endpoints, err := h.srv.GetEndpoints()
	if err != nil {
		res := utils.NewErrorResponse(500, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}
	res := utils.NewSuccessResponse(200, "all endpoints.", endpoints)
	c.JSON(res.StatusCode, res)
}

func (h *endpointHandler) CreateEndpoint(c *gin.Context) {
	var req endpoint.CreateEndpoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := utils.NewErrorResponse(400, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	endpoint, err := h.srv.CreateEndpoint(req)
	if err != nil {
		res := utils.NewErrorResponse(500, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}
	res := utils.NewSuccessResponse(200, "endpoint created.", endpoint)
	c.JSON(res.StatusCode, res)
}
