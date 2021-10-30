package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/domain/mocks"
	"github.com/apm-dev/vending-machine/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeposit(t *testing.T) {
	type args struct {
		ctx  context.Context
		coin domain.Coin
	}
	type wants struct {
		err     error
		balance uint
	}
	type testCase struct {
		name    string
		timeout time.Duration
		prepare func()
		args    args
		wants   wants
	}

	to := time.Second * 2
	ur := new(mocks.UserRepository)

	testCases := []testCase{
		{
			name:    "should succeed when buyer deposit valid coin",
			timeout: to,
			prepare: func() {
				u := &domain.User{
					Id:      1,
					Role:    domain.BUYER,
					Deposit: 0,
				}
				ur.On("FindById",
					mock.AnythingOfType("*context.timerCtx"), uint(1),
				).Return(u, nil).Once()
				ur.On("Update",
					mock.AnythingOfType("*context.timerCtx"),
					mock.AnythingOfType("*domain.User"),
				).Return(nil).Once()
			},
			args: args{
				ctx:  context.WithValue(context.Background(), domain.USER_ID_CONTEXT_KEY, uint(1)),
				coin: 50,
			},
			wants: wants{
				err:     nil,
				balance: 50,
			},
		},
	}

	for _, tc := range testCases {
		// arrange
		tc.prepare()
		svc := user.InitService(ur, nil, nil, tc.timeout)
		// action
		balance, err := svc.Deposit(tc.args.ctx, tc.args.coin)
		// assert
		if tc.wants.err != nil {
			assert.ErrorIs(t, err, tc.wants.err, tc.name)
			assert.EqualValues(t, tc.wants.balance, balance, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
			assert.EqualValues(t, tc.wants.balance, balance, tc.name)
		}
	}
	ur.AssertExpectations(t)
}
