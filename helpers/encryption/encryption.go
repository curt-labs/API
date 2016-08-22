package encryption

import (
	// "bytes"
	// "crypto/md5"
	// "errors"
	"math/rand"
	"strings"
	"time"
)

func GeneratePassword() string {
	charlist := "ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$^*?"
	charslice := strings.Split(charlist, "")
	targetlength := 8
	newpw := ""
	var index int
	for i := 0; i < targetlength; i++ {
		rand.Seed(time.Now().UnixNano())
		index = rand.Intn(len(charslice))
		newpw += charslice[index]
	}
	return newpw
}
