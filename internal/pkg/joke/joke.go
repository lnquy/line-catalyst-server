package joke

type Joker interface {
	Get() (string, error)
}
