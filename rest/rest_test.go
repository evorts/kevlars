package rest

import (
	"context"
	"encoding/json"
	"github.com/evorts/kevlars/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/reply"
	"net/http"
	"testing"
)

type TestSuite struct {
	suite.Suite
	rm   Manager
	mock *mocha.Mocha
}

type TestResponse struct {
	*contracts.ResponseSuccess[any]
	Details contracts.ErrorDetail `json:"details"`
}

func (ts *TestSuite) SetupTest() {
	ts.rm = New()
	ts.rm.SetDefaultHeaders(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	ts.mock = mocha.New(ts.T())
	ts.mock.Start()
}

func (ts *TestSuite) TestRestCall() {
	_, respBody := contracts.NewResponseOK("OK", map[string]interface{}{"status": "OK"})
	respBodyJson, _ := json.Marshal(respBody)
	scoped := ts.mock.AddMocks(
		mocha.Get(expect.URLPath("/rest-call")).
			Header("Content-Type", expect.ToEqual("application/json")).
			Header("Accept", expect.ToEqual("application/json")).
			Reply(reply.OK().Body(respBodyJson)),
	)
	var rs TestResponse
	httpCode, err := ts.rm.DisableCircuitBreaker().Get(context.Background(), ts.mock.URL()+"/rest-call", nil, &rs, nil)
	require.NoError(ts.T(), err)
	assert.Nil(ts.T(), err)
	assert.True(ts.T(), scoped.Called())
	assert.Equal(ts.T(), http.StatusOK, httpCode)
}

func (ts *TestSuite) TestCircuitBreaker() {
	ts.rm.WithCircuitBreaker(5, 5000, 1000)
	_, respBody := contracts.NewResponseOK("OK", map[string]interface{}{"status": "OK"})
	respBodyJson, _ := json.Marshal(respBody)
	scoped := ts.mock.AddMocks(
		mocha.Get(expect.URLPath("/rest-call")).
			Header("Content-Type", expect.ToEqual("application/json")).
			Header("Accept", expect.ToEqual("application/json")).
			Reply(reply.OK().Body(respBodyJson)),
	)
	var rs TestResponse
	httpCode, err := ts.rm.Get(context.Background(), ts.mock.URL()+"/rest-call", nil, &rs, nil)
	require.NoError(ts.T(), err)
	assert.Nil(ts.T(), err)
	assert.True(ts.T(), scoped.Called())
	assert.Equal(ts.T(), http.StatusOK, httpCode)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
