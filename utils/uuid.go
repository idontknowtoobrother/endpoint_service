package utils

import (
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUuid() (primitive.Binary, error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return primitive.Binary{}, err
	}
	return primitive.Binary{
		Subtype: 0x04,
		Data:    newUuid[:],
	}, nil
}

func UuidToHexString(uuidReq primitive.Binary) (string, error) {
	if uuidReq.Subtype != 0x04 || len(uuidReq.Data) != 16 {
		return "", errors.New("Invalid UUID V7 subtype")
	}
	return uuid.UUID(uuidReq.Data).String(), nil
}

func HexStringToUuid(uuidReq string) (primitive.Binary, error) {
	uuidParsed, err := uuid.Parse(uuidReq)
	if err != nil {
		return primitive.Binary{}, err
	}
	return primitive.Binary{
		Subtype: 0x04,
		Data:    uuidParsed[:],
	}, nil
}
