package validator

import (
	"testing"

	"github.com/banking/bank-server/internal/model"
)

func TestValidateStruct_Valid(t *testing.T) {
	v := New()
	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	details := v.ValidateStruct(req)
	if details != nil {
		t.Errorf("expected no validation errors, got %v", details)
	}
}

func TestValidateStruct_Invalid(t *testing.T) {
	v := New()
	req := model.LoginRequest{
		Email:    "not-an-email",
		Password: "123",
	}

	details := v.ValidateStruct(req)
	if details == nil {
		t.Fatal("expected validation errors")
	}

	if _, ok := details["email"]; !ok {
		t.Error("expected email validation error")
	}

	if _, ok := details["password"]; !ok {
		t.Error("expected password validation error")
	}
}
