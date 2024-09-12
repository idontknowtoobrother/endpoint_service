package endpoint

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EndpointRepository interface {
	GetAll() ([]Endpoint, error)
	GetByUuid(primitive.Binary) (*Endpoint, error)
	Create(Endpoint) (*Endpoint, error)
	UpdateByUuid(primitive.Binary, Endpoint) (*Endpoint, error)
	DeleteByUuid(primitive.Binary) error
}

type Endpoint struct {
	Uuid          primitive.Binary `bson:"uuid"`
	Name          string           `bson:"name"`
	Path          string           `bson:"path"`
	RedirectTo    string           `bson:"redirect_to"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at"`
	DeletedAt     time.Time        `bson:"deleted_at"`
	DeletedReason string           `bson:"deleted_reason"`
}

type endpointRepository struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewEndpointRepository(ctx context.Context, collection *mongo.Collection) EndpointRepository {
	return &endpointRepository{
		ctx:        ctx,
		collection: collection,
	}
}

func (r *endpointRepository) GetAll() ([]Endpoint, error) {
	cur, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(r.ctx)
	endpoints := make([]Endpoint, 0)
	cur.All(r.ctx, &endpoints)
	return endpoints, nil
}

func (r *endpointRepository) GetByUuid(uuid primitive.Binary) (*Endpoint, error) {
	var endpoint Endpoint
	err := r.collection.FindOne(r.ctx, bson.M{"uuid": uuid}).Decode(&endpoint)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (r *endpointRepository) Create(endpoint Endpoint) (*Endpoint, error) {
	_, err := r.collection.InsertOne(r.ctx, endpoint)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (r *endpointRepository) UpdateByUuid(uuid primitive.Binary, endpoint Endpoint) (*Endpoint, error) {
	_, err := r.collection.UpdateOne(r.ctx, bson.M{"uuid": uuid}, endpoint)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (r *endpointRepository) DeleteByUuid(uuid primitive.Binary) error {
	_, err := r.collection.DeleteOne(r.ctx, bson.M{"uuid": uuid})
	return err
}
