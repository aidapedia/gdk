package hash

type Interface interface {
	Hash(s string) string
	Verify(s string, hashed string) bool
}
