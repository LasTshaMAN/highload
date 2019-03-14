package iris_test

import (
	"highload/service/api"
	"highload/service/api/adapters/iris"
	"highload/service/api/tests"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIris(t *testing.T) {
	a := api.New()
	i := iris.New(a)
	require.NoError(t, i.Build())
	testCases := tests.NewTestSuite(i.Router)
	for name, testCase := range testCases {
		t.Run(name, testCase)
	}
}
