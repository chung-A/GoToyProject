package something

import "fmt"

func sayHello() {
	fmt.Println("say Hello")
}

func SayBye() {	
	sayHello()
	fmt.Println("say Bye!!")
}
