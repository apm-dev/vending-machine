package product

import (
	"context"
	"sync"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/algo"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

type Service struct {
	pr domain.ProductRepository
	ur domain.UserRepository
	pl sync.RWMutex
}

func InitService(pr domain.ProductRepository, ur domain.UserRepository) domain.ProductService {
	return &Service{pr: pr, ur: ur}
}

func (s *Service) Add(ctx context.Context, name string, amount uint, cost uint) (*domain.Product, error) {
	const op string = "product.service.Add"

	if cost%5 != 0 {
		return nil, domain.ErrInvalidCost
	}
	cu, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrUserNotFound
	}

	if cu.Role != domain.SELLER {
		return nil, domain.ErrPermissionDenied
	}

	p := domain.NewProduct(name, amount, cost, cu.Id)

	p.Id, err = s.pr.Insert(ctx, *p)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return p, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Product, error) {
	const op string = "product.service.List"

	s.pl.RLock()
	defer s.pl.RUnlock()

	ps, err := s.pr.List(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return ps, nil
}

func (s *Service) Update(ctx context.Context, id uint, name string, amount, cost uint) (*domain.Product, error) {
	const op string = "product.service.Update"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// only sellers can update products
	if u.Role != domain.SELLER {
		return nil, domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	p, err := s.pr.FindById(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// only related seller can update it
	if p.SellerId != u.Id {
		return nil, domain.ErrPermissionDenied
	}

	p.Name = name
	p.Count = amount
	p.Price = cost

	err = s.pr.Update(ctx, p)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	const op string = "product.service.Delete"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrUserNotFound
	}

	if u.Role != domain.SELLER {
		return domain.ErrPermissionDenied
	}

	p, err := s.pr.FindById(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	// only related seller can delete it
	if p.SellerId != u.Id {
		return domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	err = s.pr.Delete(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	return nil
}

func (s *Service) Buy(ctx context.Context, cart map[uint]uint) (*domain.Bill, error) {
	const op string = "product.service.Buy"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	if u.Role != domain.BUYER {
		return nil, domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	products := make([]domain.Product, 0, len(cart))
	items := make([]domain.Item, 0, len(products))
	var totalPrice uint

	for pid, count := range cart {
		p, err := s.pr.FindById(ctx, pid)
		if err != nil {
			logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
			return nil, domain.ErrProductNotFound
		}
		// check product availability
		if p.Count < count {
			return nil, domain.ErrInsufficientProductsAmount
		}
		// decrease product amount
		p.Count -= count
		products = append(products, *p)
		items = append(items, domain.Item{
			Name:  p.Name,
			Count: count,
			Price: count * p.Price,
		})
		// increase total price
		totalPrice += count * p.Price
		// check user balance
		if u.Deposit < totalPrice {
			return nil, domain.ErrInsufficientBalance
		}
	}

	for _, p := range products {
		//TODO: use database transaction
		err = s.pr.Update(ctx, &p)
		if err != nil {
			logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
			return nil, domain.ErrInternalServer
		}
	}

	//TODO: use database transaction
	u.Deposit -= totalPrice
	err = s.ur.Update(ctx, u)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// calculating remaining user deposit by valid coins
	refund := algo.MinimumNumberOfElementsWhoseSumIs(domain.Coins, u.Deposit)

	return &domain.Bill{
		TotalSpent: totalPrice,
		Items:      items,
		Refund:     refund,
	}, nil
}
