package id

import "github.com/jaevor/go-nanoid"

const alphabet = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

// Generate 生成一个唯一ID
func Generate() string {
	return nanoid.MustCustomASCII(alphabet, 10)()
}
