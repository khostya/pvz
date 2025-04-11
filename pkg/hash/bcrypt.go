//go:generate ${LOCAL_BIN}/ifacemaker -f ./bcrypt.go -s BcryptHash -i bcrypt -p mock_hash -o ./mocks/bcrypt.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/bcrypt.go -destination=./mocks/bcrypt.go -package=mock_hash

package hash

import "golang.org/x/crypto/bcrypt"

type BcryptHash struct {
	cost int
}
type EqualsParam struct {
	Hashed string
	V      string
}

func NewBcryptHash(cost int) *BcryptHash {
	return &BcryptHash{cost: cost}
}

func (h *BcryptHash) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}

func (h *BcryptHash) Equal(e EqualsParam) bool {
	eq := bcrypt.CompareHashAndPassword([]byte(e.Hashed), []byte(e.V))
	return eq == nil
}
