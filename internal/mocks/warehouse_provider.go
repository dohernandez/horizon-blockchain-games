// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entities "github.com/dohernandez/horizon-blockchain-games/internal/entities"

	mock "github.com/stretchr/testify/mock"
)

// WarehouseProvider is an autogenerated mock type for the WarehouseProvider type
type WarehouseProvider struct {
	mock.Mock
}

type WarehouseProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *WarehouseProvider) EXPECT() *WarehouseProvider_Expecter {
	return &WarehouseProvider_Expecter{mock: &_m.Mock}
}

// Save provides a mock function with given fields: ctx, flatten
func (_m *WarehouseProvider) Save(ctx context.Context, flatten entities.Flatten) error {
	ret := _m.Called(ctx, flatten)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Flatten) error); ok {
		r0 = rf(ctx, flatten)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WarehouseProvider_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type WarehouseProvider_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - flatten entities.Flatten
func (_e *WarehouseProvider_Expecter) Save(ctx interface{}, flatten interface{}) *WarehouseProvider_Save_Call {
	return &WarehouseProvider_Save_Call{Call: _e.mock.On("Save", ctx, flatten)}
}

func (_c *WarehouseProvider_Save_Call) Run(run func(ctx context.Context, flatten entities.Flatten)) *WarehouseProvider_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entities.Flatten))
	})
	return _c
}

func (_c *WarehouseProvider_Save_Call) Return(_a0 error) *WarehouseProvider_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WarehouseProvider_Save_Call) RunAndReturn(run func(context.Context, entities.Flatten) error) *WarehouseProvider_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewWarehouseProvider creates a new instance of WarehouseProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWarehouseProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *WarehouseProvider {
	mock := &WarehouseProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
