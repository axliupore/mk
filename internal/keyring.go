package internal

import (
	"github.com/zalando/go-keyring"
)

const service = "mk"

// Set stores a key in the system keyring.
func Set(alias, key string) error {
	return keyring.Set(service, alias, key)
}

// Get retrieves a key from the system keyring.
func Get(alias string) (string, error) {
	return keyring.Get(service, alias)
}

// Delete removes a key from the system keyring.
func Delete(alias string) error {
	return keyring.Delete(service, alias)
}

// IsNotFound reports whether the error indicates a missing key.
func IsNotFound(err error) bool {
	return err == keyring.ErrNotFound
}
