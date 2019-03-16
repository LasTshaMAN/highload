package domain_test

import (
	"highload/mocked_service/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("should return error when t2 < t1", func(t *testing.T) {
		_, err := domain.RandPointBetween(2*time.Millisecond, time.Millisecond)
		require.Error(t, err)
	})
	t.Run("should be ok when t2 == t1", func(t *testing.T) {
		t1 := time.Millisecond
		p, err := domain.RandPointBetween(t1, t1)
		require.NoError(t, err)
		require.Equal(t, p, t1)
	})
	t.Run("should be ok when t2 > t1", func(t *testing.T) {
		t1 := time.Millisecond
		t2 := 2 * time.Millisecond
		p, err := domain.RandPointBetween(t1, t2)
		require.NoError(t, err)
		require.True(t, t1 <= p && p <= t2)
	})
}
