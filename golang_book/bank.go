package main

import "fmt"

type WithD struct {
	Amount  int
	Res chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdrawCommand = make(chan WithD)
var withdrawRes = make(chan bool)

func Deposit(amount int) {
	deposits <- amount
	fmt.Println("depositing ", amount)
}
func Balance() {
	fmt.Println("balance is: ", <-balances)
}

func WithdrawFunc(amount int) bool {
	fmt.Println("trying to withdraw:", amount)
	request := WithD{amount,withdrawRes}
	withdrawCommand <- request
	if <-withdrawRes {
		fmt.Println("withdrawl was successful")
		return true
	}
	fmt.Println("withdrawl failed")
	return false
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case withdraw := <- withdrawCommand:
			if balance >= withdraw.Amount {
				balance -= withdraw.Amount
				withdraw.Res <- true
			} else {
				withdraw.Res <- false
			}
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

func main() {
	Balance()
	Deposit(100)
	Balance()
	WithdrawFunc(50)
	Balance()
	WithdrawFunc(51)
	Balance()
}
