package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	payment, principal, interest, periods, t4pe := parseFlags()
	checkFlags(payment, principal, interest, periods, t4pe)

	switch t4pe {
	case "annuity":
		computeAnniutyLoan(payment, principal, interest, periods)
	case "diff":
		computeDiffLoan(principal, interest, periods)
	}
}

func computeAnniutyLoan(payment, principal, interest float64, periods int) {
	switch {
	case payment == 0:
		p, o := calculatePayment(principal, interest, periods)
		fmt.Printf("Your monthly payment = %d!\nOverpayment = %d", p, o)
	case principal == 0:
		p, o := calculatePrincipal(payment, interest, periods)
		fmt.Printf("Your loan principal = %d!\nOverpayment = %d", p, o)
	case periods == 0:
		p, o := calculatePeriods(principal, interest, payment)
		printPeriodsResponse(p)
		fmt.Println("Overpayment =", o)
	}
}

func computeDiffLoan(principal, interest float64, periods int) {
	var total float64
	for i := float64(1); i <= float64(periods); i++ {
		d := math.Ceil(principal/float64(periods) + interest*(principal-principal*(i-1)/float64(periods)))
		fmt.Printf("Month %.0f: payment is %.0f\n", i, d)
		total += d
	}
	fmt.Println("Overpayment =", total-principal)
}

// calculatePayment calculates monthly payment base on
// p - loan principal, i - interest rate, n - number of payments
func calculatePayment(p, i float64, n int) (payment int, overpayment int) {
	payment = int(math.Ceil(p * ((i * (math.Pow(1+i, float64(n)))) / ((math.Pow(1+i, float64(n))) - 1))))
	overpayment = payment*n - int(p)
	return
}

// calculatePrincipal calculates amount of principal based on
// a - annuity payment, i - interest rate, n - number of payments
func calculatePrincipal(a, i float64, n int) (principal int, overpayment int) {
	v := math.Pow(1+i, float64(n))
	principal = int(math.Floor(a / ((i * v) / (v - 1))))
	overpayment = int((a * float64(n))) - principal
	return
}

// calculatePeriods calculates number of payments base on
// p - loan principal, i - interest rate, a - annuity payment
func calculatePeriods(p, i, a float64) (periods int, overpayment int) {
	periods = int(math.Ceil(math.Log(a/(a-i*p)) / math.Log(1+i)))
	overpayment = int(a*float64(periods) - p)
	return
}

func parseFlags() (payment, principal, interest float64, periods int, t4pe string) {
	flag.Float64Var(&payment, "payment", 0, "The payment amount")
	flag.Float64Var(&principal, "principal", 0, "Loan principal")
	flag.Float64Var(&interest, "interest", 0, "Nominal (monthly) interest rate")
	flag.IntVar(&periods, "periods", 0, "Amount of periods")
	flag.StringVar(&t4pe, "type", "unknow", "To select type of computing")
	flag.Parse()

	interest = interest / (12 * 100)
	return
}

func checkFlags(payment, principal, interest float64, periods int, t4pe string) {
	if interest == 0 || payment < 0 || principal < 0 || interest < 0 || periods < 0 {
		incorectValues()
	}

	switch params := len(os.Args); t4pe {
	case "annuity":
		if params < 5 {
			incorectValues()
		}
	case "diff":
		if payment != 0 || params < 4 {
			incorectValues()
		}
	default:
		incorectValues()
	}
}

func printPeriodsResponse(periods int) {
	switch {
	case periods <= 1:
		fmt.Println("It will take 1 month to repay this loan!")
	case periods < 12:
		fmt.Printf("It will take %d months to repay this loan!\n", periods)
	case periods >= 12:
		y := periods / 12
		m := periods % 12
		var ystr, mstr string
		if y > 1 {
			ystr = fmt.Sprintf("%d years", y)
		} else {
			ystr = "1 year"
		}

		if m == 1 {
			mstr = " and 1 month"
		} else if m > 1 {
			mstr = fmt.Sprintf(" and %d months", m)
		}
		fmt.Printf("It will take %s%s to repay this loan!\n", ystr, mstr)
	}
}

func incorectValues() {
	fmt.Println("Incorrect parameters")
	os.Exit(0)
}
