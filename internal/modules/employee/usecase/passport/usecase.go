package passport

import (
	"context"

	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository"
	"github.com/scrynaxx/tr-core/internal/modules/employee/usecase"

	"uuid"
)

type UseCase struct {
	passportRepo repository.Passport
}

func New(passportRepo repository.Passport) usecase.Passport {
	return &UseCase{
		passportRepo: passportRepo,
	}
}

func (uc *UseCase) Find(ctx context.Context, employeeID uuid.UUID) (*model.Passport, error) {
	return uc.passportRepo.Find(ctx, employeeID)
}

func (uc *UseCase) Create(ctx context.Context, employeeID uuid.UUID, input model.PassportInput) error {
	passport, err := model.NewPassport(input)
	if err != nil {
		return err
	}

	return uc.passportRepo.Create(ctx, employeeID, passport)
}

func (uc *UseCase) Update(ctx context.Context, employeeID uuid.UUID, input model.PassportInput) error {
	passport, err := uc.passportRepo.Get(ctx, employeeID)
	if err != nil {
		return err
	}

	if err = passport.Update(input); err != nil {
		return err
	}

	return uc.passportRepo.Update(ctx, employeeID, passport)
}

func (uc *UseCase) Delete(ctx context.Context, employeeID uuid.UUID) error {
	return uc.passportRepo.Delete(ctx, employeeID)
}
