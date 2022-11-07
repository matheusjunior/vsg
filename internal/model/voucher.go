package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Discount struct {
		Type  string `bson:"type" json:"type"`
		Value int    `bson:"value" json:"value"`
	}

	Voucher struct {
		Id        primitive.ObjectID `bson:"_id" json:"id"`
		UserId    uuid.UUID          `bson:"user_id" json:"user_id"`
		Username  string             `bson:"username"`
		UserEmail string             `bson:"user_email"`
		Discount  Discount           `bson:"discount" json:"discount"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		ExpireAt  time.Time          `bson:"expire_at" json:"expire_at"`
	}
)
