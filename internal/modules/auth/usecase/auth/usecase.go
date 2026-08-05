package auth

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/event"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository"
	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

type UseCase struct {
	accountRepo repository.Account
	sessionRepo repository.Session
	authorizer  security.Authorizer
	tokenizer   security.Tokenizer
	outbox      messaging.OutboxRepository
	transactor  database.Transactor
}

func New(
	accountRepo repository.Account,
	sessionRepo repository.Session,
	authorizer security.Authorizer,
	tokenizer security.Tokenizer,
	outbox messaging.OutboxRepository,
	transactor database.Transactor,
) usecase.Auth {
	return &UseCase{
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		authorizer:  authorizer,
		tokenizer:   tokenizer,
		outbox:      outbox,
		transactor:  transactor,
	}
}

func (uc *UseCase) SignIn(ctx context.Context, actor security.Actor, email, password, userAgent string) (*model.AuthState, error) {
	account, err := uc.accountRepo.GetByEmail(ctx, actor, email)
	if err != nil {
		return nil, err
	}

	if !account.PasswordEqual(password) {
		return nil, model.ErrInvalidCredentials
	}

	refresh, err := uc.tokenizer.CreateRefresh()
	if err != nil {
		return nil, err
	}

	session, err := model.NewSession(account.ID, userAgent, uc.tokenizer.HashRefresh(refresh), security.RefreshLifetime)
	if err != nil {
		return nil, err
	}

	access, err := uc.tokenizer.CreateAccess(account.ID, session.ID, actor)
	if err != nil {
		return nil, err
	}

	return database.TxResult(ctx, uc.transactor, func(ctx context.Context) (*model.AuthState, error) {
		sessionIDs, err := uc.sessionRepo.DeleteByUserAgent(ctx, account.ID, userAgent)
		if err != nil {
			return nil, err
		}

		if err = uc.storeRevoke(ctx, sessionIDs); err != nil {
			return nil, err
		}

		if err = uc.sessionRepo.Create(ctx, session); err != nil {
			return nil, err
		}

		return &model.AuthState{
			AccessToken:  access,
			RefreshToken: refresh,
			RefreshUntil: session.ExpiresAt,
			ExpiresIn:    int64(security.AccessLifetime.Seconds()),
		}, nil
	})
}

func (uc *UseCase) Refresh(ctx context.Context, actor security.Actor, refresh, userAgent string) (*model.AuthState, error) {
	refreshHash := uc.tokenizer.HashRefresh(refresh)
	session, err := uc.sessionRepo.Get(ctx, actor, refreshHash, userAgent)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		if err = database.Tx(ctx, uc.transactor, func(ctx context.Context) error {
			deleted, err := uc.sessionRepo.Delete(ctx, session.AccountID, session.ID)
			if err != nil || !deleted {
				return err
			}

			return uc.storeRevoke(ctx, []uuid.UUID{session.ID})
		}); err != nil {
			return nil, err
		}

		return nil, model.ErrSessionExpired
	}

	refresh, err = uc.tokenizer.CreateRefresh()
	if err != nil {
		return nil, err
	}

	access, err := uc.tokenizer.CreateAccess(session.AccountID, session.ID, actor)
	if err != nil {
		return nil, err
	}

	session.Rotate(uc.tokenizer.HashRefresh(refresh), security.RefreshLifetime)

	if err = uc.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return &model.AuthState{
		AccessToken:  access,
		RefreshToken: refresh,
		RefreshUntil: session.ExpiresAt,
		ExpiresIn:    int64(security.AccessLifetime.Seconds()),
	}, nil
}

func (uc *UseCase) SignOut(ctx context.Context, accountID, sessionID uuid.UUID) error {
	return database.Tx(ctx, uc.transactor, func(ctx context.Context) error {
		deleted, err := uc.sessionRepo.Delete(ctx, accountID, sessionID)
		if err != nil || !deleted {
			return err
		}

		return uc.storeRevoke(ctx, []uuid.UUID{sessionID})
	})
}

func (uc *UseCase) HandleSessionRevoked(ctx context.Context, e messaging.Event[event.SessionRevokedV1Payload]) error {
	return uc.authorizer.RevokeSessions(ctx, e.Data.SessionIDs, e.Data.RevokeUntil)
}

func (uc *UseCase) HandleEmployeeArchived(ctx context.Context, e messaging.Event[event.EmployeeArchivedV1Payload]) error {
	return database.Tx(ctx, uc.transactor, func(ctx context.Context) error {
		sessionIDs, err := uc.sessionRepo.DeleteByAccount(ctx, e.Data.AccountID)
		if err != nil {
			return err
		}

		return uc.storeRevoke(ctx, sessionIDs)
	})
}

func (uc *UseCase) storeRevoke(ctx context.Context, sessionIDs []uuid.UUID) error {
	msg, err := messaging.NewMessage(event.SessionRevokedV1, event.SessionRevokedV1Payload{
		SessionIDs: sessionIDs,
	})
	if err != nil {
		return err
	}

	return uc.outbox.StoreEvent(ctx, msg)
}
