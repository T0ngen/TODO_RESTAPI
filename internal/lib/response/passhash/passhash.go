package passhash


import (
	"crypto/md5"
	"encoding/hex"
)


func HashedPassword(password string) string{
	hash := md5.Sum([]byte(password))
	hashedPassword := hex.EncodeToString(hash[:])
	return hashedPassword
}