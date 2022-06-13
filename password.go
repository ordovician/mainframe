package mainframe

import (
	"bufio"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"os"
	"strings"
)

// Returns hashed (SHA256) version of password in base32 encoding
func HashPassword(passwd string) string {
	digest := sha256.Sum256([]byte(passwd))
	return base32.StdEncoding.EncodeToString(digest[:])
}

// Check if there is a hashed password in password file which matches
// the supplied password for given user.
// Requires that there is a passwd.txt  file with the format
//   login:passwd
// Where passwd is a SHA256 hashed password which has been
// stored in base32 encoding
func CheckLogin(user, passwd string) (bool, error) {
	file, err := os.Open("passwd.txt")

	if err != nil {
		return false, fmt.Errorf("Error when opening password file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ":", 2)

		uname := fields[0]
		if len(fields) == 2 && user == uname {
			return fields[1] == HashPassword(passwd), nil
		}
	}

	return false, nil
}
