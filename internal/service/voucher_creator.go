package service

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"matheus.com/vgs/configs"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/messaging"
	"matheus.com/vgs/internal/model"
	"matheus.com/vgs/internal/repo"
)

type voucherCreator struct {
	listener    *messaging.SQSListener
	publisher   *messaging.SQSPublisher
	voucherRepo repo.VoucherRepo
}

func NewVoucherCreatorService(listener *messaging.SQSListener, publisher *messaging.SQSPublisher, voucherRepo repo.VoucherRepo) *voucherCreator {
	return &voucherCreator{
		listener:    listener,
		publisher:   publisher,
		voucherRepo: voucherRepo,
	}
}

func (svc *voucherCreator) Start() {
	ch := svc.listener.Start()
	go func() {
		for msg := range ch {
			user := model.User{}
			if err := json.Unmarshal([]byte(msg), &user); err != nil {
				logger.Logger().Error("unmarshal user ", err)
				continue
			}
			voucher := svc.CreateVoucher(user)
			svc.PublishVoucher(voucher)
		}
	}()
}

func (svc *voucherCreator) CreateVoucher(user model.User) model.Voucher {
	now := time.Now().UTC()
	voucher := model.Voucher{
		Id:        primitive.NewObjectID(),
		UserId:    user.ID,
		UserEmail: user.Email,
		Username:  user.Name,
		Discount: model.Discount{
			Type:  "percentage",
			Value: 10,
		},
		CreatedAt: now,
		ExpireAt:  now.Add(time.Duration(configs.GetVoucherDuration()) * time.Second),
	}
	_, err := svc.voucherRepo.Insert(voucher)
	if err != nil {
		logger.Logger().Error(err)
	}
	return voucher
}

func (svc *voucherCreator) PublishVoucher(voucher model.Voucher) {
	svc.publisher.Publish(voucher)
}
