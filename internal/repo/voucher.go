package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"matheus.com/vgs/configs"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/model"
)

type VoucherRepo interface {
	Insert(voucher model.Voucher) (*mongo.InsertOneResult, error)
}

type voucherRepo struct {
	*mongo.Client
}

func NewVoucherRepo(client *mongo.Client) VoucherRepo {
	repo := voucherRepo{Client: client}
	repo.createIndex()
	return &repo
}

func (r *voucherRepo) Insert(voucher model.Voucher) (*mongo.InsertOneResult, error) {
	return r.Database("testdb").Collection("vouchers").InsertOne(context.TODO(), &voucher)
}

func (r *voucherRepo) createIndex() {
	model := mongo.IndexModel{
		Keys:    bson.D{{"created_at", 1}},
		Options: options.Index().SetExpireAfterSeconds(configs.GetVoucherDuration()),
	}
	_, err := r.Database("testdb").Collection("vouchers").Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		logger.Logger().Fatal(err)
	}
}
