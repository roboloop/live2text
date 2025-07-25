// Code generated by http://github.com/gojuno/minimock (v3.4.5). DO NOT EDIT.

package components

import (
	"context"
	"live2text/internal/services/recognition/text"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// SocketComponentMock implements SocketComponent
type SocketComponentMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcListen          func(ctx context.Context, socketPath string, formatter *text.Formatter) (err error)
	funcListenOrigin    string
	inspectFuncListen   func(ctx context.Context, socketPath string, formatter *text.Formatter)
	afterListenCounter  uint64
	beforeListenCounter uint64
	ListenMock          mSocketComponentMockListen
}

// NewSocketComponentMock returns a mock for SocketComponent
func NewSocketComponentMock(t minimock.Tester) *SocketComponentMock {
	m := &SocketComponentMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ListenMock = mSocketComponentMockListen{mock: m}
	m.ListenMock.callArgs = []*SocketComponentMockListenParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mSocketComponentMockListen struct {
	optional           bool
	mock               *SocketComponentMock
	defaultExpectation *SocketComponentMockListenExpectation
	expectations       []*SocketComponentMockListenExpectation

	callArgs []*SocketComponentMockListenParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// SocketComponentMockListenExpectation specifies expectation struct of the SocketComponent.Listen
type SocketComponentMockListenExpectation struct {
	mock               *SocketComponentMock
	params             *SocketComponentMockListenParams
	paramPtrs          *SocketComponentMockListenParamPtrs
	expectationOrigins SocketComponentMockListenExpectationOrigins
	results            *SocketComponentMockListenResults
	returnOrigin       string
	Counter            uint64
}

// SocketComponentMockListenParams contains parameters of the SocketComponent.Listen
type SocketComponentMockListenParams struct {
	ctx        context.Context
	socketPath string
	formatter  *text.Formatter
}

// SocketComponentMockListenParamPtrs contains pointers to parameters of the SocketComponent.Listen
type SocketComponentMockListenParamPtrs struct {
	ctx        *context.Context
	socketPath *string
	formatter  **text.Formatter
}

// SocketComponentMockListenResults contains results of the SocketComponent.Listen
type SocketComponentMockListenResults struct {
	err error
}

// SocketComponentMockListenOrigins contains origins of expectations of the SocketComponent.Listen
type SocketComponentMockListenExpectationOrigins struct {
	origin           string
	originCtx        string
	originSocketPath string
	originFormatter  string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmListen *mSocketComponentMockListen) Optional() *mSocketComponentMockListen {
	mmListen.optional = true
	return mmListen
}

// Expect sets up expected params for SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) Expect(ctx context.Context, socketPath string, formatter *text.Formatter) *mSocketComponentMockListen {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	if mmListen.defaultExpectation == nil {
		mmListen.defaultExpectation = &SocketComponentMockListenExpectation{}
	}

	if mmListen.defaultExpectation.paramPtrs != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by ExpectParams functions")
	}

	mmListen.defaultExpectation.params = &SocketComponentMockListenParams{ctx, socketPath, formatter}
	mmListen.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmListen.expectations {
		if minimock.Equal(e.params, mmListen.defaultExpectation.params) {
			mmListen.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmListen.defaultExpectation.params)
		}
	}

	return mmListen
}

// ExpectCtxParam1 sets up expected param ctx for SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) ExpectCtxParam1(ctx context.Context) *mSocketComponentMockListen {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	if mmListen.defaultExpectation == nil {
		mmListen.defaultExpectation = &SocketComponentMockListenExpectation{}
	}

	if mmListen.defaultExpectation.params != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Expect")
	}

	if mmListen.defaultExpectation.paramPtrs == nil {
		mmListen.defaultExpectation.paramPtrs = &SocketComponentMockListenParamPtrs{}
	}
	mmListen.defaultExpectation.paramPtrs.ctx = &ctx
	mmListen.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmListen
}

// ExpectSocketPathParam2 sets up expected param socketPath for SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) ExpectSocketPathParam2(socketPath string) *mSocketComponentMockListen {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	if mmListen.defaultExpectation == nil {
		mmListen.defaultExpectation = &SocketComponentMockListenExpectation{}
	}

	if mmListen.defaultExpectation.params != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Expect")
	}

	if mmListen.defaultExpectation.paramPtrs == nil {
		mmListen.defaultExpectation.paramPtrs = &SocketComponentMockListenParamPtrs{}
	}
	mmListen.defaultExpectation.paramPtrs.socketPath = &socketPath
	mmListen.defaultExpectation.expectationOrigins.originSocketPath = minimock.CallerInfo(1)

	return mmListen
}

// ExpectFormatterParam3 sets up expected param formatter for SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) ExpectFormatterParam3(formatter *text.Formatter) *mSocketComponentMockListen {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	if mmListen.defaultExpectation == nil {
		mmListen.defaultExpectation = &SocketComponentMockListenExpectation{}
	}

	if mmListen.defaultExpectation.params != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Expect")
	}

	if mmListen.defaultExpectation.paramPtrs == nil {
		mmListen.defaultExpectation.paramPtrs = &SocketComponentMockListenParamPtrs{}
	}
	mmListen.defaultExpectation.paramPtrs.formatter = &formatter
	mmListen.defaultExpectation.expectationOrigins.originFormatter = minimock.CallerInfo(1)

	return mmListen
}

// Inspect accepts an inspector function that has same arguments as the SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) Inspect(f func(ctx context.Context, socketPath string, formatter *text.Formatter)) *mSocketComponentMockListen {
	if mmListen.mock.inspectFuncListen != nil {
		mmListen.mock.t.Fatalf("Inspect function is already set for SocketComponentMock.Listen")
	}

	mmListen.mock.inspectFuncListen = f

	return mmListen
}

// Return sets up results that will be returned by SocketComponent.Listen
func (mmListen *mSocketComponentMockListen) Return(err error) *SocketComponentMock {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	if mmListen.defaultExpectation == nil {
		mmListen.defaultExpectation = &SocketComponentMockListenExpectation{mock: mmListen.mock}
	}
	mmListen.defaultExpectation.results = &SocketComponentMockListenResults{err}
	mmListen.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmListen.mock
}

// Set uses given function f to mock the SocketComponent.Listen method
func (mmListen *mSocketComponentMockListen) Set(f func(ctx context.Context, socketPath string, formatter *text.Formatter) (err error)) *SocketComponentMock {
	if mmListen.defaultExpectation != nil {
		mmListen.mock.t.Fatalf("Default expectation is already set for the SocketComponent.Listen method")
	}

	if len(mmListen.expectations) > 0 {
		mmListen.mock.t.Fatalf("Some expectations are already set for the SocketComponent.Listen method")
	}

	mmListen.mock.funcListen = f
	mmListen.mock.funcListenOrigin = minimock.CallerInfo(1)
	return mmListen.mock
}

// When sets expectation for the SocketComponent.Listen which will trigger the result defined by the following
// Then helper
func (mmListen *mSocketComponentMockListen) When(ctx context.Context, socketPath string, formatter *text.Formatter) *SocketComponentMockListenExpectation {
	if mmListen.mock.funcListen != nil {
		mmListen.mock.t.Fatalf("SocketComponentMock.Listen mock is already set by Set")
	}

	expectation := &SocketComponentMockListenExpectation{
		mock:               mmListen.mock,
		params:             &SocketComponentMockListenParams{ctx, socketPath, formatter},
		expectationOrigins: SocketComponentMockListenExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmListen.expectations = append(mmListen.expectations, expectation)
	return expectation
}

// Then sets up SocketComponent.Listen return parameters for the expectation previously defined by the When method
func (e *SocketComponentMockListenExpectation) Then(err error) *SocketComponentMock {
	e.results = &SocketComponentMockListenResults{err}
	return e.mock
}

// Times sets number of times SocketComponent.Listen should be invoked
func (mmListen *mSocketComponentMockListen) Times(n uint64) *mSocketComponentMockListen {
	if n == 0 {
		mmListen.mock.t.Fatalf("Times of SocketComponentMock.Listen mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmListen.expectedInvocations, n)
	mmListen.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmListen
}

func (mmListen *mSocketComponentMockListen) invocationsDone() bool {
	if len(mmListen.expectations) == 0 && mmListen.defaultExpectation == nil && mmListen.mock.funcListen == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmListen.mock.afterListenCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmListen.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Listen implements SocketComponent
func (mmListen *SocketComponentMock) Listen(ctx context.Context, socketPath string, formatter *text.Formatter) (err error) {
	mm_atomic.AddUint64(&mmListen.beforeListenCounter, 1)
	defer mm_atomic.AddUint64(&mmListen.afterListenCounter, 1)

	mmListen.t.Helper()

	if mmListen.inspectFuncListen != nil {
		mmListen.inspectFuncListen(ctx, socketPath, formatter)
	}

	mm_params := SocketComponentMockListenParams{ctx, socketPath, formatter}

	// Record call args
	mmListen.ListenMock.mutex.Lock()
	mmListen.ListenMock.callArgs = append(mmListen.ListenMock.callArgs, &mm_params)
	mmListen.ListenMock.mutex.Unlock()

	for _, e := range mmListen.ListenMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmListen.ListenMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmListen.ListenMock.defaultExpectation.Counter, 1)
		mm_want := mmListen.ListenMock.defaultExpectation.params
		mm_want_ptrs := mmListen.ListenMock.defaultExpectation.paramPtrs

		mm_got := SocketComponentMockListenParams{ctx, socketPath, formatter}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmListen.t.Errorf("SocketComponentMock.Listen got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListen.ListenMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.socketPath != nil && !minimock.Equal(*mm_want_ptrs.socketPath, mm_got.socketPath) {
				mmListen.t.Errorf("SocketComponentMock.Listen got unexpected parameter socketPath, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListen.ListenMock.defaultExpectation.expectationOrigins.originSocketPath, *mm_want_ptrs.socketPath, mm_got.socketPath, minimock.Diff(*mm_want_ptrs.socketPath, mm_got.socketPath))
			}

			if mm_want_ptrs.formatter != nil && !minimock.Equal(*mm_want_ptrs.formatter, mm_got.formatter) {
				mmListen.t.Errorf("SocketComponentMock.Listen got unexpected parameter formatter, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListen.ListenMock.defaultExpectation.expectationOrigins.originFormatter, *mm_want_ptrs.formatter, mm_got.formatter, minimock.Diff(*mm_want_ptrs.formatter, mm_got.formatter))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmListen.t.Errorf("SocketComponentMock.Listen got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmListen.ListenMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmListen.ListenMock.defaultExpectation.results
		if mm_results == nil {
			mmListen.t.Fatal("No results are set for the SocketComponentMock.Listen")
		}
		return (*mm_results).err
	}
	if mmListen.funcListen != nil {
		return mmListen.funcListen(ctx, socketPath, formatter)
	}
	mmListen.t.Fatalf("Unexpected call to SocketComponentMock.Listen. %v %v %v", ctx, socketPath, formatter)
	return
}

// ListenAfterCounter returns a count of finished SocketComponentMock.Listen invocations
func (mmListen *SocketComponentMock) ListenAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListen.afterListenCounter)
}

// ListenBeforeCounter returns a count of SocketComponentMock.Listen invocations
func (mmListen *SocketComponentMock) ListenBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListen.beforeListenCounter)
}

// Calls returns a list of arguments used in each call to SocketComponentMock.Listen.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmListen *mSocketComponentMockListen) Calls() []*SocketComponentMockListenParams {
	mmListen.mutex.RLock()

	argCopy := make([]*SocketComponentMockListenParams, len(mmListen.callArgs))
	copy(argCopy, mmListen.callArgs)

	mmListen.mutex.RUnlock()

	return argCopy
}

// MinimockListenDone returns true if the count of the Listen invocations corresponds
// the number of defined expectations
func (m *SocketComponentMock) MinimockListenDone() bool {
	if m.ListenMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.ListenMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.ListenMock.invocationsDone()
}

// MinimockListenInspect logs each unmet expectation
func (m *SocketComponentMock) MinimockListenInspect() {
	for _, e := range m.ListenMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to SocketComponentMock.Listen at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterListenCounter := mm_atomic.LoadUint64(&m.afterListenCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.ListenMock.defaultExpectation != nil && afterListenCounter < 1 {
		if m.ListenMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to SocketComponentMock.Listen at\n%s", m.ListenMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to SocketComponentMock.Listen at\n%s with params: %#v", m.ListenMock.defaultExpectation.expectationOrigins.origin, *m.ListenMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcListen != nil && afterListenCounter < 1 {
		m.t.Errorf("Expected call to SocketComponentMock.Listen at\n%s", m.funcListenOrigin)
	}

	if !m.ListenMock.invocationsDone() && afterListenCounter > 0 {
		m.t.Errorf("Expected %d calls to SocketComponentMock.Listen at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.ListenMock.expectedInvocations), m.ListenMock.expectedInvocationsOrigin, afterListenCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *SocketComponentMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockListenInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *SocketComponentMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *SocketComponentMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockListenDone()
}
