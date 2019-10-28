// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
	"time"

	"code.cloudfoundry.org/clock"
)

type FakeClock struct {
	AfterStub        func(time.Duration) <-chan time.Time
	afterMutex       sync.RWMutex
	afterArgsForCall []struct {
		arg1 time.Duration
	}
	afterReturns struct {
		result1 <-chan time.Time
	}
	afterReturnsOnCall map[int]struct {
		result1 <-chan time.Time
	}
	NewTickerStub        func(time.Duration) clock.Ticker
	newTickerMutex       sync.RWMutex
	newTickerArgsForCall []struct {
		arg1 time.Duration
	}
	newTickerReturns struct {
		result1 clock.Ticker
	}
	newTickerReturnsOnCall map[int]struct {
		result1 clock.Ticker
	}
	NewTimerStub        func(time.Duration) clock.Timer
	newTimerMutex       sync.RWMutex
	newTimerArgsForCall []struct {
		arg1 time.Duration
	}
	newTimerReturns struct {
		result1 clock.Timer
	}
	newTimerReturnsOnCall map[int]struct {
		result1 clock.Timer
	}
	NowStub        func() time.Time
	nowMutex       sync.RWMutex
	nowArgsForCall []struct {
	}
	nowReturns struct {
		result1 time.Time
	}
	nowReturnsOnCall map[int]struct {
		result1 time.Time
	}
	SinceStub        func(time.Time) time.Duration
	sinceMutex       sync.RWMutex
	sinceArgsForCall []struct {
		arg1 time.Time
	}
	sinceReturns struct {
		result1 time.Duration
	}
	sinceReturnsOnCall map[int]struct {
		result1 time.Duration
	}
	SleepStub        func(time.Duration)
	sleepMutex       sync.RWMutex
	sleepArgsForCall []struct {
		arg1 time.Duration
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClock) After(arg1 time.Duration) <-chan time.Time {
	fake.afterMutex.Lock()
	ret, specificReturn := fake.afterReturnsOnCall[len(fake.afterArgsForCall)]
	fake.afterArgsForCall = append(fake.afterArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("After", []interface{}{arg1})
	fake.afterMutex.Unlock()
	if fake.AfterStub != nil {
		return fake.AfterStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.afterReturns
	return fakeReturns.result1
}

func (fake *FakeClock) AfterCallCount() int {
	fake.afterMutex.RLock()
	defer fake.afterMutex.RUnlock()
	return len(fake.afterArgsForCall)
}

func (fake *FakeClock) AfterCalls(stub func(time.Duration) <-chan time.Time) {
	fake.afterMutex.Lock()
	defer fake.afterMutex.Unlock()
	fake.AfterStub = stub
}

func (fake *FakeClock) AfterArgsForCall(i int) time.Duration {
	fake.afterMutex.RLock()
	defer fake.afterMutex.RUnlock()
	argsForCall := fake.afterArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClock) AfterReturns(result1 <-chan time.Time) {
	fake.afterMutex.Lock()
	defer fake.afterMutex.Unlock()
	fake.AfterStub = nil
	fake.afterReturns = struct {
		result1 <-chan time.Time
	}{result1}
}

func (fake *FakeClock) AfterReturnsOnCall(i int, result1 <-chan time.Time) {
	fake.afterMutex.Lock()
	defer fake.afterMutex.Unlock()
	fake.AfterStub = nil
	if fake.afterReturnsOnCall == nil {
		fake.afterReturnsOnCall = make(map[int]struct {
			result1 <-chan time.Time
		})
	}
	fake.afterReturnsOnCall[i] = struct {
		result1 <-chan time.Time
	}{result1}
}

func (fake *FakeClock) NewTicker(arg1 time.Duration) clock.Ticker {
	fake.newTickerMutex.Lock()
	ret, specificReturn := fake.newTickerReturnsOnCall[len(fake.newTickerArgsForCall)]
	fake.newTickerArgsForCall = append(fake.newTickerArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("NewTicker", []interface{}{arg1})
	fake.newTickerMutex.Unlock()
	if fake.NewTickerStub != nil {
		return fake.NewTickerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newTickerReturns
	return fakeReturns.result1
}

func (fake *FakeClock) NewTickerCallCount() int {
	fake.newTickerMutex.RLock()
	defer fake.newTickerMutex.RUnlock()
	return len(fake.newTickerArgsForCall)
}

func (fake *FakeClock) NewTickerCalls(stub func(time.Duration) clock.Ticker) {
	fake.newTickerMutex.Lock()
	defer fake.newTickerMutex.Unlock()
	fake.NewTickerStub = stub
}

func (fake *FakeClock) NewTickerArgsForCall(i int) time.Duration {
	fake.newTickerMutex.RLock()
	defer fake.newTickerMutex.RUnlock()
	argsForCall := fake.newTickerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClock) NewTickerReturns(result1 clock.Ticker) {
	fake.newTickerMutex.Lock()
	defer fake.newTickerMutex.Unlock()
	fake.NewTickerStub = nil
	fake.newTickerReturns = struct {
		result1 clock.Ticker
	}{result1}
}

func (fake *FakeClock) NewTickerReturnsOnCall(i int, result1 clock.Ticker) {
	fake.newTickerMutex.Lock()
	defer fake.newTickerMutex.Unlock()
	fake.NewTickerStub = nil
	if fake.newTickerReturnsOnCall == nil {
		fake.newTickerReturnsOnCall = make(map[int]struct {
			result1 clock.Ticker
		})
	}
	fake.newTickerReturnsOnCall[i] = struct {
		result1 clock.Ticker
	}{result1}
}

func (fake *FakeClock) NewTimer(arg1 time.Duration) clock.Timer {
	fake.newTimerMutex.Lock()
	ret, specificReturn := fake.newTimerReturnsOnCall[len(fake.newTimerArgsForCall)]
	fake.newTimerArgsForCall = append(fake.newTimerArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("NewTimer", []interface{}{arg1})
	fake.newTimerMutex.Unlock()
	if fake.NewTimerStub != nil {
		return fake.NewTimerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newTimerReturns
	return fakeReturns.result1
}

func (fake *FakeClock) NewTimerCallCount() int {
	fake.newTimerMutex.RLock()
	defer fake.newTimerMutex.RUnlock()
	return len(fake.newTimerArgsForCall)
}

func (fake *FakeClock) NewTimerCalls(stub func(time.Duration) clock.Timer) {
	fake.newTimerMutex.Lock()
	defer fake.newTimerMutex.Unlock()
	fake.NewTimerStub = stub
}

func (fake *FakeClock) NewTimerArgsForCall(i int) time.Duration {
	fake.newTimerMutex.RLock()
	defer fake.newTimerMutex.RUnlock()
	argsForCall := fake.newTimerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClock) NewTimerReturns(result1 clock.Timer) {
	fake.newTimerMutex.Lock()
	defer fake.newTimerMutex.Unlock()
	fake.NewTimerStub = nil
	fake.newTimerReturns = struct {
		result1 clock.Timer
	}{result1}
}

func (fake *FakeClock) NewTimerReturnsOnCall(i int, result1 clock.Timer) {
	fake.newTimerMutex.Lock()
	defer fake.newTimerMutex.Unlock()
	fake.NewTimerStub = nil
	if fake.newTimerReturnsOnCall == nil {
		fake.newTimerReturnsOnCall = make(map[int]struct {
			result1 clock.Timer
		})
	}
	fake.newTimerReturnsOnCall[i] = struct {
		result1 clock.Timer
	}{result1}
}

func (fake *FakeClock) Now() time.Time {
	fake.nowMutex.Lock()
	ret, specificReturn := fake.nowReturnsOnCall[len(fake.nowArgsForCall)]
	fake.nowArgsForCall = append(fake.nowArgsForCall, struct {
	}{})
	fake.recordInvocation("Now", []interface{}{})
	fake.nowMutex.Unlock()
	if fake.NowStub != nil {
		return fake.NowStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nowReturns
	return fakeReturns.result1
}

func (fake *FakeClock) NowCallCount() int {
	fake.nowMutex.RLock()
	defer fake.nowMutex.RUnlock()
	return len(fake.nowArgsForCall)
}

func (fake *FakeClock) NowCalls(stub func() time.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = stub
}

func (fake *FakeClock) NowReturns(result1 time.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	fake.nowReturns = struct {
		result1 time.Time
	}{result1}
}

func (fake *FakeClock) NowReturnsOnCall(i int, result1 time.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	if fake.nowReturnsOnCall == nil {
		fake.nowReturnsOnCall = make(map[int]struct {
			result1 time.Time
		})
	}
	fake.nowReturnsOnCall[i] = struct {
		result1 time.Time
	}{result1}
}

func (fake *FakeClock) Since(arg1 time.Time) time.Duration {
	fake.sinceMutex.Lock()
	ret, specificReturn := fake.sinceReturnsOnCall[len(fake.sinceArgsForCall)]
	fake.sinceArgsForCall = append(fake.sinceArgsForCall, struct {
		arg1 time.Time
	}{arg1})
	fake.recordInvocation("Since", []interface{}{arg1})
	fake.sinceMutex.Unlock()
	if fake.SinceStub != nil {
		return fake.SinceStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sinceReturns
	return fakeReturns.result1
}

func (fake *FakeClock) SinceCallCount() int {
	fake.sinceMutex.RLock()
	defer fake.sinceMutex.RUnlock()
	return len(fake.sinceArgsForCall)
}

func (fake *FakeClock) SinceCalls(stub func(time.Time) time.Duration) {
	fake.sinceMutex.Lock()
	defer fake.sinceMutex.Unlock()
	fake.SinceStub = stub
}

func (fake *FakeClock) SinceArgsForCall(i int) time.Time {
	fake.sinceMutex.RLock()
	defer fake.sinceMutex.RUnlock()
	argsForCall := fake.sinceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClock) SinceReturns(result1 time.Duration) {
	fake.sinceMutex.Lock()
	defer fake.sinceMutex.Unlock()
	fake.SinceStub = nil
	fake.sinceReturns = struct {
		result1 time.Duration
	}{result1}
}

func (fake *FakeClock) SinceReturnsOnCall(i int, result1 time.Duration) {
	fake.sinceMutex.Lock()
	defer fake.sinceMutex.Unlock()
	fake.SinceStub = nil
	if fake.sinceReturnsOnCall == nil {
		fake.sinceReturnsOnCall = make(map[int]struct {
			result1 time.Duration
		})
	}
	fake.sinceReturnsOnCall[i] = struct {
		result1 time.Duration
	}{result1}
}

func (fake *FakeClock) Sleep(arg1 time.Duration) {
	fake.sleepMutex.Lock()
	fake.sleepArgsForCall = append(fake.sleepArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("Sleep", []interface{}{arg1})
	fake.sleepMutex.Unlock()
	if fake.SleepStub != nil {
		fake.SleepStub(arg1)
	}
}

func (fake *FakeClock) SleepCallCount() int {
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	return len(fake.sleepArgsForCall)
}

func (fake *FakeClock) SleepCalls(stub func(time.Duration)) {
	fake.sleepMutex.Lock()
	defer fake.sleepMutex.Unlock()
	fake.SleepStub = stub
}

func (fake *FakeClock) SleepArgsForCall(i int) time.Duration {
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	argsForCall := fake.sleepArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClock) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.afterMutex.RLock()
	defer fake.afterMutex.RUnlock()
	fake.newTickerMutex.RLock()
	defer fake.newTickerMutex.RUnlock()
	fake.newTimerMutex.RLock()
	defer fake.newTimerMutex.RUnlock()
	fake.nowMutex.RLock()
	defer fake.nowMutex.RUnlock()
	fake.sinceMutex.RLock()
	defer fake.sinceMutex.RUnlock()
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClock) recordInvocation(key string, args []interface{}) {
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

var _ clock.Clock = new(FakeClock)