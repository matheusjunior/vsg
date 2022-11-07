package service

import (
	"encoding/json"

	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/messaging"
	"matheus.com/vgs/internal/model"
	"matheus.com/vgs/internal/sender"
)

type notifier struct {
	listener *messaging.SQSListener
	sender   sender.Sender
}

func NewNotifierService(listener *messaging.SQSListener, sender sender.Sender) *notifier {
	return &notifier{
		listener: listener,
		sender:   sender,
	}
}

func (svc *notifier) Start() {
	ch := svc.listener.Start()
	go func() {
		for msg := range ch {
			voucher := model.Voucher{}
			if err := json.Unmarshal([]byte(msg), &voucher); err != nil {
				logger.Logger().Error("unmarshal voucher", err)
				continue
			}
			logger.Logger().Infof("notify voucher %+v\n", voucher)
			svc.sender.Send(voucher)
		}
	}()
}
