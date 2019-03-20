// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import gqlschema "github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"

import mock "github.com/stretchr/testify/mock"
import v1beta1 "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"

// gqlInstanceConverter is an autogenerated mock type for the gqlInstanceConverter type
type gqlInstanceConverter struct {
	mock.Mock
}

// ToGQL provides a mock function with given fields: in
func (_m *gqlInstanceConverter) ToGQL(in *v1beta1.ServiceInstance) (*gqlschema.ServiceInstance, error) {
	ret := _m.Called(in)

	var r0 *gqlschema.ServiceInstance
	if rf, ok := ret.Get(0).(func(*v1beta1.ServiceInstance) *gqlschema.ServiceInstance); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gqlschema.ServiceInstance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1beta1.ServiceInstance) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
