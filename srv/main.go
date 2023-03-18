package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println(isVietnamesePhoneNumber("0388888888"))
}
func isVietnamesePhoneNumber(number string) bool {
	// return ``.test(number)
	ok, _ := regexp.MatchString("/^(0?)(3[2-9]|5[6|8|9]|7[0|6-9]|8[0-6|8|9]|9[0-4|6-9])[0-9]{7}$/", number)
	return ok
}
