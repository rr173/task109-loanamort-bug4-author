package loan_test

import (
	"context"
	"testing"

	"task109-loanamort/internal/loan"
	"task109-loanamort/internal/store"
)

func TestAsOfIncludesPaymentAtRequestedPeriod(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(t.TempDir() + "/loan.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := loan.New(db)
	b, err := svc.CreateBorrower(ctx, loan.CreateBorrowerRequest{Name: "borrower"})
	if err != nil {
		t.Fatal(err)
	}
	l, err := svc.CreateLoan(ctx, loan.CreateLoanRequest{BorrowerID: b.BorrowerID, PrincipalCents: 100000, AnnualPercent: 12, Periods: 4, Type: loan.EqualInstallment})
	if err != nil {
		t.Fatal(err)
	}
	s, err := svc.Schedule(ctx, l.LoanID)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = svc.RecordPayment(ctx, l.LoanID, loan.RecordPaymentRequest{AmountCents: s.Periods[0].Payment}); err != nil {
		t.Fatal(err)
	}
	payments, err := svc.ListPaymentsByLoan(ctx, l.LoanID)
	if err != nil {
		t.Fatal(err)
	}
	if len(payments) != 1 || payments[0].Seq != 1 {
		t.Fatalf("first payment sequence=%d, want 1", payments[0].Seq)
	}
	bal, err := svc.Balance(ctx, l.LoanID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if bal.AsOf != 1 {
		t.Fatalf("as_of=%d, want 1", bal.AsOf)
	}
	want := l.Principal - s.Periods[0].Principal
	if bal.Outstanding != want {
		t.Fatalf("outstanding=%d, want %d", bal.Outstanding, want)
	}
}
