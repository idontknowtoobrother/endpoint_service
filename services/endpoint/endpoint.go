package endpoint

import (
	"context"
	"errors"
	"fmt"

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
}

type EndpointResponse struct {
	Uuid          string `json:"uuid"`
	Name          string `bson:"name"`
	Path          string `json:"path"`
	RedirectTo    string `json:"redirect_to"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
	DeletedReason string `json:"deleted_reason"`
}

type CreateEndpoinRequest struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	RedirectTo string `json:"redirect_to",binding:"required"`
}

type endpointService struct {
	ctx  context.Context
	repo endpoint.EndpointRepository
}

func NewEndpointService(ctx context.Context, repo endpoint.EndpointRepository) EndpointService {
	return &endpointService{
		ctx:  ctx,
		repo: repo,
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
		Path:          req.Path,
		RedirectTo:    req.RedirectTo,
		CreatedAt:     utils.TimeNow(),
		UpdatedAt:     utils.TimeNow(),
		DeletedAt:     utils.TimeZero(),
		DeletedReason: "",
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
		DeletedAt:     utils.TimeToString(newEndpoint.DeletedAt),
		DeletedReason: newEndpoint.DeletedReason,
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

		endpointResponses = append(endpointResponses, EndpointResponse{
			Uuid:          uuidHexString,
			Name:          endpoint.Name,
			Path:          endpoint.Path,
			RedirectTo:    endpoint.RedirectTo,
			CreatedAt:     utils.TimeToString(endpoint.CreatedAt),
			UpdatedAt:     utils.TimeToString(endpoint.UpdatedAt),
			DeletedAt:     utils.TimeToString(endpoint.DeletedAt),
			DeletedReason: endpoint.DeletedReason,
		})
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

	if endpoint.DeletedReason != "" && !endpoint.DeletedAt.IsZero() {
		return nil, errors.New(fmt.Sprintf(constraints.ErrEndpointHasBeenDeleted, endpoint.DeletedReason))
	}

	uuidHexString, err := utils.UuidToHexString(endpoint.Uuid)
	if err != nil {
		return nil, err
	}

	return &EndpointResponse{
		Uuid:          uuidHexString,
		Name:          endpoint.Name,
		Path:          endpoint.Path,
		RedirectTo:    endpoint.RedirectTo,
		CreatedAt:     utils.TimeToString(endpoint.CreatedAt),
		UpdatedAt:     utils.TimeToString(endpoint.UpdatedAt),
		DeletedAt:     utils.TimeToString(endpoint.DeletedAt),
		DeletedReason: endpoint.DeletedReason,
	}, nil
}
