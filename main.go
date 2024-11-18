package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"strings"
)

type Loan struct {
	payment     float64
	principal   float64
	periods     int
	interest    float64
	paymentType string
}

func main() {
	loan, err := NewLoan()

	if err != nil {
		fmt.Println(err)
	} else {
		calculate(loan)
	}
}

func NewLoan() (*Loan, error) {
	loan := Loan{}

	flag.Float64Var(&loan.interest, "interest", 0, "a charge for borrowing ")
	flag.Float64Var(&loan.payment, "payment", 0, "payment amount is precisely this fixed sum of money that you need to pay at regular intervals.")
	flag.IntVar(&loan.periods, "periods", 0, "the number of months in which repayments will be made")
	flag.Float64Var(&loan.principal, "principal", 0, "loan principal")
	flag.StringVar(&loan.paymentType, "type", "", "the type of payment")

	flag.Parse()

	if len(flag.Args()) > 3 || loan.isValid() {
		return &loan, nil
	}

	return nil, errors.New("Incorrect parameters")
}

func calculate(loan *Loan) {
	switch loan.paymentType {
	case "diff":
		calculateDifferentiatedPayment(loan)
	case "annuity":
		switch {
		case loan.principal > 0 && loan.payment > 0:
			calculateNumberOfPayments(loan)
		case loan.principal > 0 && loan.periods > 0:
			calculateAnnuityPayment(loan)
		case loan.payment > 0 && loan.periods > 0:
			calculateLoanPrincipal(loan)
		}
	}
}

func calculateNumberOfPayments(loan *Loan) {
	i := loan.interestValue()

	n := math.Log(loan.payment/(loan.payment-i*loan.principal)) / math.Log(1+i)

	months := int(math.Ceil(n)) // Round up to the next whole number
	years := months / 12
	remainingMonths := months % 12

	switch {
	case years == 0:
		fmt.Printf("It will take %d %s to repay this loan!\n", remainingMonths, periodInPlural("month", years))
	case remainingMonths == 0:
		fmt.Printf("It will take %d %s to repay this loan!\n", years, periodInPlural("year", years))
	default:
		fmt.Printf("It will take %d %s and %d %s to repay this loan!\n",
			years, periodInPlural("year", years), remainingMonths, periodInPlural("month", years))
	}

	totalPayment := float64(months * int(loan.payment))
	fmt.Printf("Overpayment = %.0f\n", totalPayment-loan.principal)
}

func periodInPlural(period string, n int) string {
	if n > 1 {
		return period + "s"
	}
	return period
}

func calculateAnnuityPayment(loan *Loan) {
	i := loan.interestValue()
	annuityPayment := loan.principal * i * math.Pow(1+i, float64(loan.periods)) / (math.Pow(1+i, float64(loan.periods)) - 1)

	annuityPayment = math.Ceil(annuityPayment)

	fmt.Printf("Your monthly payment = %.0f!\n", annuityPayment)
	totalPayment := annuityPayment * float64(loan.periods)
	fmt.Printf("Overpayment = %.0f\n", totalPayment-loan.principal)
}

func calculateDifferentiatedPayment(loan *Loan) {
	i := loan.interestValue()

	totalPayment := float64(0)

	for m := 1; m <= loan.periods; m++ {
		differentiatedPayment := loan.principal/float64(loan.periods) + i*(loan.principal-loan.principal*float64(m-1)/float64(loan.periods))
		differentiatedPayment = math.Ceil(differentiatedPayment)
		totalPayment += differentiatedPayment
		fmt.Printf("Month %d: payment is %.0f\n", m, differentiatedPayment)
	}
	fmt.Printf("Overpayment = %.0f\n", totalPayment-loan.principal)
}

func calculateLoanPrincipal(loan *Loan) {
	i := loan.interestValue()
	principal := loan.payment / (i * math.Pow(1+i, float64(loan.periods)) / (math.Pow(1+i, float64(loan.periods)) - 1))
	fmt.Printf("Your loan principal = %.0f!\n", principal)
}

func (l *Loan) interestValue() float64 {
	return l.interest / (12 * 100) // Convert annual interest rate to monthly and to a decimal
}

func (l *Loan) isValid() bool {
	switch {
	case !strings.EqualFold(l.paymentType, "annuity") && !strings.EqualFold(l.paymentType, "diff"):
		return false
	case l.payment < 0 || l.interest < 0 || l.periods < 0 || l.principal < 0:
		return false
	case strings.EqualFold(l.paymentType, "diff") && l.payment > 0:
		return false
	case strings.EqualFold(l.paymentType, "annuity") && l.interest <= 0:
		return false
	default:
		return true
	}
}
