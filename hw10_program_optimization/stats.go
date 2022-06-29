package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var errUnmarshallingJSON = errors.New("unmarshalling json error")

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

//go:generate easyjson -all stats.go

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		user := User{}
		err := user.UnmarshalJSON(line)
		if err != nil {
			return nil, errUnmarshallingJSON
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			split := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[split]
			num++
			result[split] = num
		}
	}

	return result, nil
}
