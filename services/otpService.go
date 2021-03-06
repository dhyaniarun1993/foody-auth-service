package services

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type otpService struct {
	otpRepository repositories.OtpRepository
}

// NewOtpService creates and return otp service object
func NewOtpService(otpRepository repositories.OtpRepository) OtpService {
	return &otpService{otpRepository}
}

func (service *otpService) Validate(ctx context.Context, phoneNumber string,
	otp string, client models.Client) (bool, errors.AppError) {

	key := "otpLogin:" + client.ID + ":" + phoneNumber
	otpFromRepository, getOtpErr := service.otpRepository.Get(ctx, key)
	if getOtpErr != nil {
		return false, getOtpErr
	}

	if otpFromRepository == "" {
		return false, nil
	}

	if otp != otpFromRepository {
		return false, nil
	}

	deleteOtpErr := service.otpRepository.Delete(ctx, key)
	if deleteOtpErr != nil {
		return false, deleteOtpErr
	}

	return true, nil
}

func (service *otpService) Generate(ctx context.Context, phoneNumber string,
	client models.Client) errors.AppError {

	// Todo: Instead tp using phone number, use uuid which will be sent back to the client
	// Verification will also be done using the uuid
	key := "otpLogin:" + client.ID + ":" + phoneNumber
	minOtp := 100000
	maxOtp := 999999
	rand.Seed(time.Now().Unix())
	otp := rand.Intn(maxOtp-minOtp) + minOtp

	setOtpError := service.otpRepository.Set(ctx, key, strconv.Itoa(otp))
	if setOtpError != nil {
		return setOtpError
	}

	// Todo: Sms otp to the user
	return nil
}
