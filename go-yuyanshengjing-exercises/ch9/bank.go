package ch9

type withdraw struct {
	amount int
	res    chan bool
}

var deposits = make(chan int)        // send amount to deposit
var balances = make(chan int)        // receive balance
var withdraws = make(chan *withdraw) // send amount to withdraw

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	wd := &withdraw{
		amount: amount,
		res:    make(chan bool),
	}
	withdraws <- wd

	return <-wd.res
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case amount := <-withdraws:
			if balance-amount.amount < 0 {
				amount.res <- false
				break
			}
			balance -= amount.amount
			amount.res <- true
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
