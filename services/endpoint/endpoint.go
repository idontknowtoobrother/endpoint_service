package endpoint

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/idontknowtoobrother/practice_go_hexagonal/constraints"
	"github.com/idontknowtoobrother/practice_go_hexagonal/repositories/endpoint"
	"github.com/idontknowtoobrother/practice_go_hexagonal/utils"
	"github.com/xyproto/randomstring"
	"go.mongodb.org/mongo-driver/mongo"
)

type EndpointService interface {
	GetEndpoints() ([]EndpointResponse, error)
	GetEndpoint(uuid string) (*EndpointResponse, error)
	CreateEndpoint(req CreateEndpoinRequest) (*EndpointResponse, error)
	RedirectToDiscordWebhook(path string, body map[string]interface{}) (*resty.Response, error)
	UpdateEndpoint(req UpdateEndpointRequest) (*EndpointResponse, error)
	DeleteEndpoint(req DeleteEndpointRequest) (*EndpointResponse, error)
}

type EndpointResponse struct {
	Uuid          string  `json:"uuid"`
	Name          string  `bson:"name"`
	Path          string  `json:"path"`
	RedirectTo    string  `json:"redirect_to"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at,omitempty"`
	DeletedReason *string `json:"deleted_reason,omitempty"`
}

type CreateEndpoinRequest struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	RedirectTo string `json:"redirect_to" binding:"required"`
}

type UpdateEndpointRequest struct {
	Uuid       string  `json:"uuid" binding:"required"`
	Name       *string `json:"name,omitempty"`
	RedirectTo *string `json:"redirect_to,omitempty"`
}

type DeleteEndpointRequest struct {
	Uuid   string `json:"uuid" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type endpointService struct {
	ctx         context.Context
	repo        endpoint.EndpointRepository
	restyClient *resty.Client
}

func NewEndpointService(ctx context.Context, repo endpoint.EndpointRepository, restyClient *resty.Client) EndpointService {
	return &endpointService{
		ctx:         ctx,
		repo:        repo,
		restyClient: restyClient,
	}
}

func (s *endpointService) CreateEndpoint(req CreateEndpoinRequest) (*EndpointResponse, error) {
	newUuid, err := utils.NewUuid()
	if err != nil {
		return nil, err
	}

	newPath := randomstring.CookieFriendlyString(16)
	_, err = s.repo.GetByPath(newPath)
	for err != mongo.ErrNoDocuments {
		newPath = randomstring.CookieFriendlyString(16)
	}

	endpoint := endpoint.Endpoint{
		Uuid:          newUuid,
		Name:          req.Name,
		Path:          newPath,
		RedirectTo:    req.RedirectTo,
		CreatedAt:     utils.TimeNow(),
		UpdatedAt:     utils.TimeNow(),
		DeletedAt:     nil,
		DeletedReason: nil,
	}

	newEndpoint, err := s.repo.Create(endpoint)
	if err != nil {
		return nil, err
	}

	uuidHexString, err := utils.UuidToHexString(newEndpoint.Uuid)
	if err != nil {
		return nil, err
	}

	return &EndpointResponse{
		Uuid:          uuidHexString,
		Name:          newEndpoint.Name,
		Path:          newEndpoint.Path,
		RedirectTo:    newEndpoint.RedirectTo,
		CreatedAt:     utils.TimeToString(newEndpoint.CreatedAt),
		UpdatedAt:     utils.TimeToString(newEndpoint.UpdatedAt),
		DeletedAt:     nil,
		DeletedReason: nil,
	}, nil
}

func (s *endpointService) GetEndpoints() ([]EndpointResponse, error) {
	endpoints, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	endpointResponses := make([]EndpointResponse, 0)
	for _, endpoint := range endpoints {
		uuidHexString, err := utils.UuidToHexString(endpoint.Uuid)
		if err != nil {
			return nil, err
		}
		resEndpoint := EndpointResponse{
			Uuid:       uuidHexString,
			Name:       endpoint.Name,
			Path:       endpoint.Path,
			RedirectTo: endpoint.RedirectTo,
			CreatedAt:  utils.TimeToString(endpoint.CreatedAt),
			UpdatedAt:  utils.TimeToString(endpoint.UpdatedAt),
		}
		if endpoint.DeletedAt != nil {
			deletedAt := utils.TimeToString(*endpoint.DeletedAt)
			resEndpoint.DeletedAt = &deletedAt
			resEndpoint.DeletedReason = endpoint.DeletedReason
		}
		endpointResponses = append(endpointResponses, resEndpoint)
	}

	return endpointResponses, nil
}

func (s *endpointService) GetEndpoint(uuid string) (*EndpointResponse, error) {
	uuidBinary, err := utils.HexStringToUuid(uuid)
	if err != nil {
		return nil, err
	}

	endpoint, err := s.repo.GetByUuid(uuidBinary)
	if err != nil {
		return nil, err
	}

	if endpoint.DeletedReason != nil && endpoint.DeletedAt != nil {
		return nil, errors.New(fmt.Sprintf(constraints.ErrEndpointHasBeenDeleted, *endpoint.DeletedReason))
	}

	uuidHexString, err := utils.UuidToHexString(endpoint.Uuid)
	if err != nil {
		return nil, err
	}
	endpointResponse := &EndpointResponse{
		Uuid:       uuidHexString,
		Name:       endpoint.Name,
		Path:       endpoint.Path,
		RedirectTo: endpoint.RedirectTo,
		CreatedAt:  utils.TimeToString(endpoint.CreatedAt),
		UpdatedAt:  utils.TimeToString(endpoint.UpdatedAt),
	}

	if endpoint.DeletedAt != nil {
		deletedAt := utils.TimeToString(*endpoint.DeletedAt)
		endpointResponse.DeletedAt = &deletedAt
		endpointResponse.DeletedReason = endpoint.DeletedReason
	}

	return endpointResponse, nil
}

func (s *endpointService) RedirectToDiscordWebhook(path string, body map[string]interface{}) (*resty.Response, error) {

	endpoint, err := s.repo.GetByPath(path)
	if err != nil {
		return nil, err
	}

	if endpoint.DeletedReason != nil && endpoint.DeletedAt != nil {
		return nil, errors.New(fmt.Sprintf(constraints.ErrEndpointHasBeenDeleted, *endpoint.DeletedReason))
	}

	req := s.restyClient.R().SetBody(body)
	response, err := req.Post(endpoint.RedirectTo)
	if response.StatusCode() != 200 {
		return response, err
	}
	return response, nil
}

func (s *endpointService) UpdateEndpoint(req UpdateEndpointRequest) (*EndpointResponse, error) {
	uuidBinary, err := utils.HexStringToUuid(req.Uuid)
	if err != nil {
		return nil, err
	}

	endpoint, err := s.repo.GetByUuid(uuidBinary)
	if err != nil {
		return nil, err
	}

	if endpoint.DeletedReason != nil && endpoint.DeletedAt != nil {
		return nil, errors.New(fmt.Sprintf(constraints.ErrEndpointHasBeenDeleted, *endpoint.DeletedReason))
	}

	if req.Name != nil {
		endpoint.Name = *req.Name
	}
	if req.RedirectTo != nil {
		endpoint.RedirectTo = *req.RedirectTo
	}

	endpoint.UpdatedAt = utils.TimeNow()

	updatedEndpoint, err := s.repo.UpdateByUuid(uuidBinary, endpoint)
	if err != nil {
		return nil, err
	}

	uuidHexString, err := utils.UuidToHexString(updatedEndpoint.Uuid)
	if err != nil {
		return nil, err
	}

	return &EndpointResponse{
		Uuid:          uuidHexString,
		Name:          updatedEndpoint.Name,
		Path:          updatedEndpoint.Path,
		RedirectTo:    updatedEndpoint.RedirectTo,
		CreatedAt:     utils.TimeToString(updatedEndpoint.CreatedAt),
		UpdatedAt:     utils.TimeToString(updatedEndpoint.UpdatedAt),
		DeletedAt:     nil,
		DeletedReason: nil,
	}, nil
}

func (s *endpointService) DeleteEndpoint(req DeleteEndpointRequest) (*EndpointResponse, error) {
	uuidBinary, err := utils.HexStringToUuid(req.Uuid)
	if err != nil {
		return nil, err
	}

	endpoint, err := s.repo.GetByUuid(uuidBinary)
	if err != nil {
		return nil, err
	}

	if endpoint.DeletedReason != nil && endpoint.DeletedAt != nil {
		return nil, errors.New(fmt.Sprintf(constraints.ErrEndpointHasBeenDeleted, *endpoint.DeletedReason))
	}

	deletedAt := utils.TimeNow()
	deletedAtString := utils.TimeToString(deletedAt)
	endpoint.DeletedAt = &deletedAt
	endpoint.DeletedReason = &req.Reason

	deletedEndpoint, err := s.repo.UpdateByUuid(uuidBinary, endpoint)
	if err != nil {
		return nil, err
	}

	uuidHexString, err := utils.UuidToHexString(deletedEndpoint.Uuid)
	if err != nil {
		return nil, err
	}

	return &EndpointResponse{
		Uuid:          uuidHexString,
		Name:          deletedEndpoint.Name,
		Path:          deletedEndpoint.Path,
		RedirectTo:    deletedEndpoint.RedirectTo,
		CreatedAt:     utils.TimeToString(deletedEndpoint.CreatedAt),
		UpdatedAt:     utils.TimeToString(deletedEndpoint.UpdatedAt),
		DeletedAt:     &deletedAtString,
		DeletedReason: deletedEndpoint.DeletedReason,
	}, nil
}
