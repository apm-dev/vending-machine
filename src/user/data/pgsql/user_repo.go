package pgsql

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func InitUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Insert(ctx context.Context, u domain.User) (uint, error) {
	const op string = "user.data.pgsql.user_repo.Insert"

	dbUser := new(User)
	dbUser.FromDomain(&u)

	err := r.db.WithContext(ctx).Create(&dbUser).Error
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	return dbUser.ID, nil
}

func (r *UserRepository) FindById(ctx context.Context, id uint) (*domain.User, error) {
	const op string = "user.data.pgsql.user_repo.FindById"

	dbUser := new(User)

	err := r.db.WithContext(ctx).First(&dbUser, "id = ?", id).Error
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return dbUser.ToDomain(), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, un string) (*domain.User, error) {
	const op string = "user.data.pgsql.user_repo.FindByUsername"

	dbUser := new(User)

	err := r.db.WithContext(ctx).First(&dbUser, "username = ?", un).Error
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return dbUser.ToDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	const op string = "user.data.pgsql.user_repo.Update"

	dbUser := new(User)
	dbUser.FromDomain(u)

	err := r.db.WithContext(ctx).Where("id = ?", u.Id).Updates(&dbUser).Error
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	const op string = "user.data.pgsql.user_repo.Delete"

	err := r.db.WithContext(ctx).Delete(&User{}, id).Error
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
