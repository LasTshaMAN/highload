package iris_test

import (
	"highload/service/api"
	"highload/service/api/adapters/iris"
	"highload/service/api/tests"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIris(t *testing.T) {
	testCases := tests.NewTestSuite()
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testCase(t, irisBootstrap)
		})
	}
}

func irisBootstrap(t *testing.T, a *api.API) http.Handler {
	i := iris.New(a)
	require.NoError(t, i.Build())
	return i.Router
}
