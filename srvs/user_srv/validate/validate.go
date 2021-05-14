package validate

import (
	"errors"
	"regexp"
	"unicode"
)

func hasSpace(str string) bool {
	for _, letter := range str {
		if unicode.IsSpace(letter) {
			return true
		}
	}
	return false
}

// IsUsername 验证用户名合法性，用户名不能为空，长度为5-20
func IsUsername(username string) error {
	if hasSpace(username) {
		return errors.New("用户名不合法")
	} else if len(username) < 5 || len(username) > 20 {
		return errors.New("用户名长度需为5-20")
	}
	//matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	//if err != nil || !matched {
	//	return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	//}
	//matched, err = regexp.MatchString("^[a-zA-Z]", username)
	//if err != nil || !matched {
	//	return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	//}
	return nil
}

// IsEmail 验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if hasSpace(email) {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}


// IsPassword 是否是合法的密码
func IsPassword(password, rePassword string) error {
	if hasSpace(password) {
		return errors.New("密码不符合规范")
	}
	pattern := `^[a-zA-Z0-9]{7,30}$`
	if matched, _ := regexp.MatchString(pattern, password); !matched {
		return errors.New("密码不符合规范")
	}
	if password != rePassword {
		return errors.New("两次输入密码不匹配")
	}
	return nil
}