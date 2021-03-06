// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
)

type DefaultACLProvider struct {
	CheckACLStub        func(string, string, interface{}) error
	checkACLMutex       sync.RWMutex
	checkACLArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 interface{}
	}
	checkACLReturns struct {
		result1 error
	}
	checkACLReturnsOnCall map[int]struct {
		result1 error
	}
	CheckACLNoChannelStub        func(string, interface{}) error
	checkACLNoChannelMutex       sync.RWMutex
	checkACLNoChannelArgsForCall []struct {
		arg1 string
		arg2 interface{}
	}
	checkACLNoChannelReturns struct {
		result1 error
	}
	checkACLNoChannelReturnsOnCall map[int]struct {
		result1 error
	}
	IsPtypePolicyStub        func(string) bool
	isPtypePolicyMutex       sync.RWMutex
	isPtypePolicyArgsForCall []struct {
		arg1 string
	}
	isPtypePolicyReturns struct {
		result1 bool
	}
	isPtypePolicyReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DefaultACLProvider) CheckACL(arg1 string, arg2 string, arg3 interface{}) error {
	fake.checkACLMutex.Lock()
	ret, specificReturn := fake.checkACLReturnsOnCall[len(fake.checkACLArgsForCall)]
	fake.checkACLArgsForCall = append(fake.checkACLArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 interface{}
	}{arg1, arg2, arg3})
	stub := fake.CheckACLStub
	fakeReturns := fake.checkACLReturns
	fake.recordInvocation("CheckACL", []interface{}{arg1, arg2, arg3})
	fake.checkACLMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DefaultACLProvider) CheckACLCallCount() int {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	return len(fake.checkACLArgsForCall)
}

func (fake *DefaultACLProvider) CheckACLCalls(stub func(string, string, interface{}) error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = stub
}

func (fake *DefaultACLProvider) CheckACLArgsForCall(i int) (string, string, interface{}) {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	argsForCall := fake.checkACLArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *DefaultACLProvider) CheckACLReturns(result1 error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = nil
	fake.checkACLReturns = struct {
		result1 error
	}{result1}
}

func (fake *DefaultACLProvider) CheckACLReturnsOnCall(i int, result1 error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = nil
	if fake.checkACLReturnsOnCall == nil {
		fake.checkACLReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkACLReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DefaultACLProvider) CheckACLNoChannel(arg1 string, arg2 interface{}) error {
	fake.checkACLNoChannelMutex.Lock()
	ret, specificReturn := fake.checkACLNoChannelReturnsOnCall[len(fake.checkACLNoChannelArgsForCall)]
	fake.checkACLNoChannelArgsForCall = append(fake.checkACLNoChannelArgsForCall, struct {
		arg1 string
		arg2 interface{}
	}{arg1, arg2})
	stub := fake.CheckACLNoChannelStub
	fakeReturns := fake.checkACLNoChannelReturns
	fake.recordInvocation("CheckACLNoChannel", []interface{}{arg1, arg2})
	fake.checkACLNoChannelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DefaultACLProvider) CheckACLNoChannelCallCount() int {
	fake.checkACLNoChannelMutex.RLock()
	defer fake.checkACLNoChannelMutex.RUnlock()
	return len(fake.checkACLNoChannelArgsForCall)
}

func (fake *DefaultACLProvider) CheckACLNoChannelCalls(stub func(string, interface{}) error) {
	fake.checkACLNoChannelMutex.Lock()
	defer fake.checkACLNoChannelMutex.Unlock()
	fake.CheckACLNoChannelStub = stub
}

func (fake *DefaultACLProvider) CheckACLNoChannelArgsForCall(i int) (string, interface{}) {
	fake.checkACLNoChannelMutex.RLock()
	defer fake.checkACLNoChannelMutex.RUnlock()
	argsForCall := fake.checkACLNoChannelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DefaultACLProvider) CheckACLNoChannelReturns(result1 error) {
	fake.checkACLNoChannelMutex.Lock()
	defer fake.checkACLNoChannelMutex.Unlock()
	fake.CheckACLNoChannelStub = nil
	fake.checkACLNoChannelReturns = struct {
		result1 error
	}{result1}
}

func (fake *DefaultACLProvider) CheckACLNoChannelReturnsOnCall(i int, result1 error) {
	fake.checkACLNoChannelMutex.Lock()
	defer fake.checkACLNoChannelMutex.Unlock()
	fake.CheckACLNoChannelStub = nil
	if fake.checkACLNoChannelReturnsOnCall == nil {
		fake.checkACLNoChannelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkACLNoChannelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DefaultACLProvider) IsPtypePolicy(arg1 string) bool {
	fake.isPtypePolicyMutex.Lock()
	ret, specificReturn := fake.isPtypePolicyReturnsOnCall[len(fake.isPtypePolicyArgsForCall)]
	fake.isPtypePolicyArgsForCall = append(fake.isPtypePolicyArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.IsPtypePolicyStub
	fakeReturns := fake.isPtypePolicyReturns
	fake.recordInvocation("IsPtypePolicy", []interface{}{arg1})
	fake.isPtypePolicyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DefaultACLProvider) IsPtypePolicyCallCount() int {
	fake.isPtypePolicyMutex.RLock()
	defer fake.isPtypePolicyMutex.RUnlock()
	return len(fake.isPtypePolicyArgsForCall)
}

func (fake *DefaultACLProvider) IsPtypePolicyCalls(stub func(string) bool) {
	fake.isPtypePolicyMutex.Lock()
	defer fake.isPtypePolicyMutex.Unlock()
	fake.IsPtypePolicyStub = stub
}

func (fake *DefaultACLProvider) IsPtypePolicyArgsForCall(i int) string {
	fake.isPtypePolicyMutex.RLock()
	defer fake.isPtypePolicyMutex.RUnlock()
	argsForCall := fake.isPtypePolicyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DefaultACLProvider) IsPtypePolicyReturns(result1 bool) {
	fake.isPtypePolicyMutex.Lock()
	defer fake.isPtypePolicyMutex.Unlock()
	fake.IsPtypePolicyStub = nil
	fake.isPtypePolicyReturns = struct {
		result1 bool
	}{result1}
}

func (fake *DefaultACLProvider) IsPtypePolicyReturnsOnCall(i int, result1 bool) {
	fake.isPtypePolicyMutex.Lock()
	defer fake.isPtypePolicyMutex.Unlock()
	fake.IsPtypePolicyStub = nil
	if fake.isPtypePolicyReturnsOnCall == nil {
		fake.isPtypePolicyReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isPtypePolicyReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *DefaultACLProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	fake.checkACLNoChannelMutex.RLock()
	defer fake.checkACLNoChannelMutex.RUnlock()
	fake.isPtypePolicyMutex.RLock()
	defer fake.isPtypePolicyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DefaultACLProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
