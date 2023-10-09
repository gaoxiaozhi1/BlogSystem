package pwd

import (
	"fmt"
	"testing"
)

func TestHashPwd(t *testing.T) {
	fmt.Println(HashPwd("123456"))
}

func TestCheckPwd(t *testing.T) {
	fmt.Println(CheckPwd("$2a$04$euuEw/be3ldyl7Xbwklmne5Lr1KSaWcjV3i8YPHw5sZNkUOZdJ2yS", "123456"))
}
