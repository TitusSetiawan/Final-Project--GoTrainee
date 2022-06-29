package svr

import (
	"net/mail"
)

func ValidationUsersEmail(email string) (string, bool) {
	if email == "" {
		return "", false
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return email, true
}

func ValidationUsersPass(pass string) bool {
	if pass == "" {
		return false
	} else {
		return len(pass) >= 6
	}
}

func ValidationUsersUname(uname string) bool {
	return uname != ""
}

func ValidationUsersAge(age int) bool {
	if age == 0 {
		return false
	}
	return age >= 8
}
