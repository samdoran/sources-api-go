package middleware

import (
	"net/http"
	"strings"
	"testing"

	"github.com/RedHatInsights/sources-api-go/internal/testutils/request"
	"github.com/RedHatInsights/sources-api-go/internal/testutils/templates"
	"github.com/labstack/echo/v4"
)

// idValidationFunc is a helper variable which makes the context to actually handle the errors produced by the
// middleware.
var idValidationFunc = HandleErrors(
	IdValidation(
		func(c echo.Context) error { return nil },
	),
)

// uuidValidationFunc is a helper variable which makes the context to actually handle the errors produced by the
// middleware.
var uuidValidationFunc = HandleErrors(
	UuidValidation(
		func(c echo.Context) error { return nil },
	),
)

// TestExtractValidateId tests that the function under test works as expected when a valid parameter name and ID are
// provided.
func TestExtractValidateId(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)

	paramName := "id"
	paramValue := "12345"

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := idValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a valid ID: %s`, err)
	}

	want := http.StatusOK
	got := rec.Code

	if want != got {
		t.Errorf(`unexpected status code received. Want "%d", got "%d"`, want, got)
	}
}

// TestExtractValidateIdNonParseable tests that a bad request response is returned when a non-parseable ID is set.
func TestExtractValidateIdNonParseable(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)

	paramName := "id"
	paramValue := "abcde"

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := idValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a non parseable ID: %s`, err)
	}

	templates.BadRequestTest(t, rec)

	want := "could not parse the provided ID"
	got := rec.Body.String()

	if !strings.Contains(got, want) {
		t.Errorf(`unexpected error received when testing for a non parseable ID. Want an error containing "%s", got "%s"`, want, got)
	}
}

// TestExtractValidateIdGreaterZero tests that a bad request response is returned when an ID which is not greater than
// zero is set.
func TestExtractValidateIdGreaterZero(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)

	paramName := "id"
	paramValue := "0"

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := idValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a not greater than zero ID: %s`, err)
	}

	templates.BadRequestTest(t, rec)

	want := "the provided ID must be greater than zero"
	got := rec.Body.String()

	if !strings.Contains(got, want) {
		t.Errorf(`unexpected error received when testing for a non greater than zero ID. Want "%s", got "%s"`, want, got)
	}
}

// TestExtractValidateAuthUuid tests that the function under test properly works when the UUIDs are set on the right
// parameter name.
func TestExtractValidateAuthUuid(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)

	paramName := "uid"
	paramValue := "15"

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := uuidValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a valid ID: %s`, err)
	}

	want := http.StatusOK
	got := rec.Code

	if want != got {
		t.Errorf(`unexpected status code received. Want "%d", got "%d"`, want, got)
	}
}

// TestExtractValidateAuthUuidWrongParamName tests that a bad request is returned when the UUID has been set on a parameter
// name different to the one expected.
func TestExtractValidateAuthUuidWrongParamName(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)
	originalSecretStore := conf.SecretStore
	conf.SecretStore = "vault"

	paramName := "uuid"
	paramValue := ""

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := uuidValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a valid ID: %s`, err)
	}

	templates.BadRequestTest(t, rec)

	want := "the UUID cannot be empty or missing"
	got := rec.Body.String()

	if !strings.Contains(got, want) {
		t.Errorf(`unexpected error when validating an UUID which was incorrectly set in a different parameter name. Want "%s", got "%s"`, want, got)
	}

	conf.SecretStore = originalSecretStore
}

// TestExtractValidateAuthUuidEmpty tests that an error is returned when the given UUID is empty.
func TestExtractValidateAuthUuidEmpty(t *testing.T) {
	c, rec := request.CreateTestContext(http.MethodGet, "/", nil, nil)
	originalSecretStore := conf.SecretStore
	conf.SecretStore = "vault"

	paramName := "uid"
	paramValue := ""

	c.SetParamNames(paramName)
	c.SetParamValues(paramValue)

	err := uuidValidationFunc(c)
	if err != nil {
		t.Errorf(`unexpected error received when a validating a valid ID: %s`, err)
	}

	templates.BadRequestTest(t, rec)

	want := "the UUID cannot be empty or missing"
	got := rec.Body.String()

	if !strings.Contains(got, want) {
		t.Errorf(`unexpected error when validating an empty UUID. Want "%s", got "%s"`, want, got)
	}

	conf.SecretStore = originalSecretStore
}
