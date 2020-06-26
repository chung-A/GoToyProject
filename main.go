package main

//1. 패키지 명은 가장 끝단의 폴더명과 동일하게 해야함.

import (
	"fmt"

	"github.com/chung-A/GoToyProject/Wow"
	"github.com/chung-A/GoToyProject/something"
)

func main() {
	fmt.Println("메인시작")
	something.SayBye()
	something.HelloWorld()
	Wow.Wow()
}
