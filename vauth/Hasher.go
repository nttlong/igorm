package unvsauth

type Hasher interface {
	Hash(password string) (string, error)
	Verify(hashed, password string) bool
}
