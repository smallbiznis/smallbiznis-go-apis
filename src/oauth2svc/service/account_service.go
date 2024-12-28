package service

import (
	"context"
	"fmt"

	"github.com/smallbiznis/oauth2-server/internal/pkg/errors"
	"github.com/smallbiznis/oauth2-server/internal/pkg/strings"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
)

type IAccountService interface {
	HandleSignUp(context.Context, model.RequestSignUp) (*model.AggregateAccount, error)
	HandleSignInWithPassword(context.Context, model.RequestSignInWithPassword) (*model.AggregateAccount, error)
	HandleSignInWithPhoneNumber(context.Context, model.RequestSignInWithPhoneNumber) (*model.AggregateAccount, error)
	HandleSendVerificationCode(context.Context, model.RequestSendVerificationCode) (*model.AggregateAccount, error)
}

type accountService struct {
	applicationRepository repository.IApplicationRepository
	accountRepository     repository.IAccountRepository
	sessionRepository     repository.ISessionRepository
}

func NewAccountService(
	applicationRepository repository.IApplicationRepository,
	accountRepository repository.IAccountRepository,
	sessionRepository repository.ISessionRepository,
) IAccountService {
	return &accountService{
		applicationRepository,
		accountRepository,
		sessionRepository,
	}
}

func (s *accountService) HandleSignUp(ctx context.Context, params model.RequestSignUp) (*model.AggregateAccount, error) {
	tenant := ctx.Value("tenant").(*model.Organization)
	account := model.Account{
		OrganizationID: tenant.ID,
		Type:           model.User,
		Provider:       params.Provider,
	}

	if params.Provider == model.Password {
		account.Username = params.Email
	} else if params.Provider == model.PhoneNumber {
		account.Username = params.PhoneNumber
	}

	exist, err := s.accountRepository.FindOne(ctx, account)
	if err != nil {
		return nil, err
	}

	if exist != nil {
		return nil, errors.ErrInvalidEmailOrPassord
	}

	account.Profile = model.Profile{
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}

	if params.Provider == model.Password {
		account.SetPassword(params.Password)
		account.Profile.ProfileURI = fmt.Sprintf("https://gravatar.com/avatar/%s", strings.Hash256(account.Username))
	}

	account.Roles = append(account.Roles, string(model.ROLE_USER))
	return s.accountRepository.Save(ctx, account)
}

func (s *accountService) HandleSignInWithPassword(ctx context.Context, params model.RequestSignInWithPassword) (*model.AggregateAccount, error) {
	tenant := ctx.Value("tenant").(*model.Organization)
	exist, err := s.accountRepository.FindOne(ctx, model.Account{
		OrganizationID: tenant.ID,
		Provider:       model.Password,
		Username:       params.Email,
	})
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.ErrInvalidEmailOrPassord
	}

	if !exist.ComparePassword(params.Password) {
		return nil, errors.ErrInvalidEmailOrPassord
	}

	newSess := model.NewSession(params.Request)
	newSess.ID = strings.RandomHex(32)
	newSess.UserID = exist.ID
	if _, err := s.sessionRepository.Create(ctx, *newSess); err != nil {
		return nil, err
	}

	exist.SessionID = newSess.ID
	return exist, nil
}

func (s *accountService) HandleSignInWithPhoneNumber(ctx context.Context, params model.RequestSignInWithPhoneNumber) (*model.AggregateAccount, error) {
	return nil, errors.ErrNotImplement
	// exist, err := s.accountRepository.FindOne(ctx, model.Account{
	// 	Provider:    model.PhoneNumber,
	// 	PhoneNumber: params.PhoneNumber,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// if exist == nil {
	// 	return nil, errors.ErrInvalidEmailOrPassord
	// }

	// code, err := s.verificationRepository.FindOne(ctx, model.VerificationCode{
	// 	ID:        params.SessionID,
	// 	AccountID: exist.ID,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// if code == nil {
	// 	return nil, errors.ErrInvalidSessionID
	// }

	// if code.Revoke {
	// 	return nil, errors.ErrInvalidVerificationCode
	// }

	// if code.ExpiredAt.Compare(time.Now()) > -1 {
	// 	return nil, errors.ErrVerificationCodeExpired
	// }

	// return nil, nil
}

func (s *accountService) HandleSendVerificationCode(ctx context.Context, params model.RequestSendVerificationCode) (*model.AggregateAccount, error) {
	tenant := ctx.Value("tenant").(*model.Organization)
	exist, err := s.accountRepository.FindOne(ctx, model.Account{
		OrganizationID: tenant.ID,
		Provider:       model.PhoneNumber,
		Username:       params.PhoneNumber,
	})
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.ErrInvalidEmailOrPassord
	}

	return nil, nil
}
