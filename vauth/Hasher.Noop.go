package vauth

type NoopHasher struct{}

func (n *NoopHasher) Hash(password string) (string, error) {
	return password, nil
}

func (n *NoopHasher) Verify(hashed, password string) bool {
	return hashed == password
}
