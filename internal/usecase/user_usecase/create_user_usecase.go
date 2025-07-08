package user_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"

	"github.com/google/uuid"
)

func (u *UserUseCase) CreateUser(
	ctx context.Context, user UserInputDTO) (*UserOutputDTO, *internal_error.InternalError) {

	userEntity := user_entity.User{
		Id:   uuid.New().String(),
		Name: user.Name,
	}

	err := u.UserRepository.CreateUser(
		ctx, userEntity)
	if err != nil {
		return &UserOutputDTO{}, err
	}

	return &UserOutputDTO{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}, nil
}
