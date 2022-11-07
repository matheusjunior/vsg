package configs

const (
	userQueueUrl    = "http://localhost:4566/000000000000/USER"
	voucherQueueUrl = "http://localhost:4566/000000000000/VOUCHER"
)

func GetUserQueueUrl() string {
	return userQueueUrl
}

func GetVoucherQueueUrl() string {
	return voucherQueueUrl
}
