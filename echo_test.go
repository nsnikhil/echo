package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	host    = "127.0.0.1:8085"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type echoTestSuite struct {
	suite.Suite
}

func (ets *echoTestSuite) SetupSuite() {
	go startServer(host)
	time.Sleep(time.Second)
}

func TestEchoSuite(t *testing.T) {
	suite.Run(t, new(echoTestSuite))
}

func (ets *echoTestSuite) TestEchoSuccess() {
	for i := 0; i < 10; i++ {
		val := randString()
		resp, err := http.Get(buildURL(ets.T(), map[string]string{"value": val}))
		ets.Require().NoError(err)

		ets.Require().Equal(http.StatusOK, resp.StatusCode)

		b, err := ioutil.ReadAll(resp.Body)
		ets.Require().NoError(err)

		ets.Assert().Equal(val, string(b))
	}

}

func (ets *echoTestSuite) TestEchoFailure() {
	testCases := map[string]map[string]string{
		"test failure when value is missing": {"value": ""},
		"test failure when value is empty":   {},
	}

	for name, params := range testCases {
		ets.Run(name, func() {
			resp, err := http.Get(buildURL(ets.T(), params))
			ets.Require().NoError(err)

			ets.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func buildURL(t *testing.T, params map[string]string) string {
	u, err := url.Parse(fmt.Sprintf("http://%s/echo", host))
	require.NoError(t, err)

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func randString() string {
	randLetters := func(n int) string {
		rand.Seed(time.Now().UnixNano())

		sz := len(letters)

		sb := strings.Builder{}
		sb.Grow(n)

		for i := 0; i < n; i++ {
			sb.WriteByte(letters[rand.Intn(sz)])
		}

		return sb.String()
	}

	randInt := func(min, max int) int {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(max-min+1) + min
	}

	return randLetters(randInt(8, 14))
}
