package main

import "fmt"

/*삽질 기록*/
//1. 패키지 명은 가장 끝단의 폴더명과 동일하게 해야함.

func main() {
	println(deferFunc("HeyHey defer test"))
}

//...으로 하면 동일한 타입의 여러 매개변수를 받을 수 있다.
//_하면 무시가능.
func arrayArgFunc(len int, txt ...string) int {
	fmt.Println("myFunc")
	fmt.Println(txt)
	return 0
}

//naked 기능.
func nakedFunc(name string) (length int, orignTxt string) {
	length = len(name)
	orignTxt = name
	return
}

//defer기능-함수가 끝나고 원래자리로 가기전에 defer내용을 실행시키고 감.
//defer 내부에 함수안에 매개변수로 함수가 있으면 미리 실행시키고 가는듯...?
//println(함수()) 이런식으로는 안쓰는게 좋을 듯 하다.
func deferFunc(testStr string) (rst string) {
	defer println(len(testStr))
	println(testStr + "/ deferFunc")
	return testStr
}

//loop
