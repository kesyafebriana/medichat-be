package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/util"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type paymentService struct {
	dataRepository domain.DataRepository
	cloudProvider  util.CloudinaryProvider
}

type PaymentServiceOpts struct {
	DataRepository domain.DataRepository
	CloudProvider  util.CloudinaryProvider
}

func NewPaymentService(opts PaymentServiceOpts) *paymentService {
	return &paymentService{
		dataRepository: opts.DataRepository,
		cloudProvider:  opts.CloudProvider,
	}
}

func (s *paymentService) List(
	ctx context.Context,
	dets domain.PaymentListDetails,
) ([]domain.Payment, domain.PageInfo, error) {
	paymentRepo := s.dataRepository.PaymentRepository()
	accountRepo := s.dataRepository.AccountRepository()
	userRepo := s.dataRepository.UserRepository()

	accountID, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	account, err := accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	if account.Role == domain.AccountRoleUser {
		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return nil, domain.PageInfo{}, apperror.Wrap(err)
		}

		dets.UserID = &user.ID
	}

	page, err := paymentRepo.GetPageInfo(ctx, dets)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	payments, err := paymentRepo.List(ctx, dets)
	if err != nil {
		return nil, domain.PageInfo{}, apperror.Wrap(err)
	}

	return payments, page, err
}

func (s *paymentService) GetByID(
	ctx context.Context,
	id int64,
) (domain.Payment, error) {
	paymentRepo := s.dataRepository.PaymentRepository()

	payment, err := paymentRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Payment{}, apperror.Wrap(err)
	}

	return payment, err
}

func (s *paymentService) GetByInvoiceNumber(
	ctx context.Context,
	num string,
) (domain.Payment, error) {
	paymentRepo := s.dataRepository.PaymentRepository()
	accountRepo := s.dataRepository.AccountRepository()
	userRepo := s.dataRepository.UserRepository()

	accountID, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return domain.Payment{}, apperror.Wrap(err)
	}

	account, err := accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return domain.Payment{}, apperror.Wrap(err)
	}

	payment, err := paymentRepo.GetByInvoiceNumber(ctx, num)
	if err != nil {
		return domain.Payment{}, apperror.Wrap(err)
	}

	if account.Role == domain.AccountRoleUser {
		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return domain.Payment{}, apperror.Wrap(err)
		}

		if payment.User.ID != user.ID {
			return domain.Payment{}, apperror.NewForbidden(nil)
		}
	}

	return payment, err
}

func (s *paymentService) UploadPaymentClosure(
	ctx context.Context,
	num string,
	file multipart.File,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		paymentRepo := dr.PaymentRepository()
		userRepo := dr.UserRepository()
		orderRepo := dr.OrderRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		user, err := userRepo.GetByAccountID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		payment, err := paymentRepo.GetByInvoiceNumber(ctx, num)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if payment.User.ID != user.ID {
			return nil, apperror.NewForbidden(nil)
		}
		if payment.FileURL != nil {
			return nil, apperror.NewPaymentAlreadyUploaded(nil)
		}

		if file == nil {
			return nil, apperror.NewInternalFmt("file should not be nil")
		}

		res, err := s.cloudProvider.UploadImage(ctx, file, uploader.UploadParams{})
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		payment.FileURL = &res.SecureURL

		payment, err = paymentRepo.Update(ctx, payment)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = orderRepo.UpdateStatusByPaymentID(
			ctx, payment.ID,
			domain.OrderStatusWaitingConfirmation,
		)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *paymentService) UploadPayment(
	ctx context.Context,
	num string,
	file multipart.File,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.UploadPaymentClosure(ctx, num, file),
	)
	return err
}

func (s *paymentService) ConfirmPaymentClosure(
	ctx context.Context,
	num string,
) domain.AtomicFunc[any] {
	return func(dr domain.DataRepository) (any, error) {
		paymentRepo := dr.PaymentRepository()
		accuntRepo := dr.AccountRepository()
		orderRepo := dr.OrderRepository()

		accountID, err := util.GetAccountIDFromContext(ctx)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		account, err := accuntRepo.GetByID(ctx, accountID)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if account.Role != domain.AccountRoleAdmin {
			return nil, apperror.NewForbidden(nil)
		}

		payment, err := paymentRepo.GetByInvoiceNumber(ctx, num)
		if err != nil {
			return nil, apperror.Wrap(err)
		}
		if payment.FileURL == nil {
			return nil, apperror.NewPaymentNotYetUploaded(nil)
		}
		if payment.IsConfirmed {
			return nil, apperror.NewPaymentAlreadyConfirmed(nil)
		}

		payment.IsConfirmed = true

		payment, err = paymentRepo.Update(ctx, payment)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		err = orderRepo.UpdateStatusByPaymentID(
			ctx, payment.ID,
			domain.OrderStatusProcessing,
		)
		if err != nil {
			return nil, apperror.Wrap(err)
		}

		return nil, nil
	}
}

func (s *paymentService) ConfirmPayment(
	ctx context.Context,
	num string,
) error {
	_, err := domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.ConfirmPaymentClosure(ctx, num),
	)
	return err
}
