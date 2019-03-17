//go:generate mockgen -source avg.go -destination avg.mock.gen.go -package domain

package domain

type Avg interface {
	Value() (int, error)
}
