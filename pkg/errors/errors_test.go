package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNotFound(t *testing.T) {
	err := NotFound("not here")
	if err.Code != http.StatusNotFound {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusNotFound)
	}
	if err.Message != "not here" {
		t.Errorf("Message = %q, want %q", err.Message, "not here")
	}
	if err.Error() != "not here" {
		t.Errorf("Error() = %q, want %q", err.Error(), "not here")
	}
}

func TestNotFoundErr(t *testing.T) {
	inner := errors.New("db error")
	err := NotFoundErr("not found", inner)
	if err.Code != http.StatusNotFound {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusNotFound)
	}
	if !errors.Is(err, inner) {
		t.Error("errors.Is should find inner error")
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("bad")
	if err.Code != http.StatusBadRequest {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusBadRequest)
	}
}

func TestBadRequestErr(t *testing.T) {
	inner := errors.New("parse error")
	err := BadRequestErr("bad request", inner)
	if !errors.Is(err, inner) {
		t.Error("errors.Is should find inner error")
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("no auth")
	if err.Code != http.StatusUnauthorized {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusUnauthorized)
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("forbidden")
	if err.Code != http.StatusForbidden {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusForbidden)
	}
}

func TestConflict(t *testing.T) {
	err := Conflict("conflict")
	if err.Code != http.StatusConflict {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusConflict)
	}
}

func TestInternal(t *testing.T) {
	err := Internal("oops")
	if err.Code != http.StatusInternalServerError {
		t.Errorf("Code = %d, want %d", err.Code, http.StatusInternalServerError)
	}
}

func TestInternalErr(t *testing.T) {
	inner := errors.New("internal")
	err := InternalErr("oops", inner)
	if !errors.Is(err, inner) {
		t.Error("errors.Is should find inner error")
	}
}

func TestFromError_AppError(t *testing.T) {
	original := NotFound("already app error")
	result := FromError(original)
	if result != original {
		t.Error("FromError should return AppError as-is")
	}
}

func TestFromError_RegularError(t *testing.T) {
	result := FromError(errors.New("generic"))
	if result.Code != http.StatusInternalServerError {
		t.Errorf("Code = %d, want %d", result.Code, http.StatusInternalServerError)
	}
}

func TestUnwrap(t *testing.T) {
	inner := errors.New("root cause")
	wrapper := InternalErr("wrapped", inner)
	if !errors.Is(wrapper, inner) {
		t.Error("Unwrap should allow errors.Is to find inner")
	}
}
