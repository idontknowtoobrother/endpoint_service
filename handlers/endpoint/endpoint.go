package endpoint

import (
	"encoding/json"

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

func (h *endpointHandler) SendToWebhook(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		res := utils.NewErrorResponse(400, constraints.ErrNotFound)
		c.JSON(res.StatusCode, res)
		return
	}

	body := c.Request.Body
	defer body.Close()
	bodyInterface := map[string]interface{}{}
	if err := c.ShouldBindJSON(&bodyInterface); err != nil {
		res := utils.NewErrorResponse(400, "invalid body")
		c.JSON(res.StatusCode, res)
		return
	}

	result, err := h.srv.RedirectToDiscordWebhook(path, bodyInterface)
	if err != nil {
		res := utils.NewErrorResponse(500, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	responseFromDiscord := map[string]interface{}{}
	json.Unmarshal(result.Body(), &responseFromDiscord)
	c.JSON(result.StatusCode(), responseFromDiscord)
}

func (h *endpointHandler) UpdateEndpoint(c *gin.Context) {
	var req endpoint.UpdateEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := utils.NewErrorResponse(400, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	endpoint, err := h.srv.UpdateEndpoint(req)
	if err != nil {
		res := utils.NewErrorResponse(500, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	res := utils.NewSuccessResponse(200, "endpoint updated.", endpoint)
	c.JSON(res.StatusCode, res)
}

func (h *endpointHandler) DeleteEndpoint(c *gin.Context) {
	var req endpoint.DeleteEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := utils.NewErrorResponse(400, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	result, err := h.srv.DeleteEndpoint(req)
	if err != nil {
		res := utils.NewErrorResponse(500, err.Error())
		c.JSON(res.StatusCode, res)
		return
	}

	res := utils.NewSuccessResponse(200, "endpoint deleted.", result)
	c.JSON(res.StatusCode, res)
}
