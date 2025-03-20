package test_jwt

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	jwt "github.com/Manolo-Esc/gommence/src/internal/infra/jwt"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
	"github.com/Manolo-Esc/gommence/src/tests/libtest"
	"github.com/stretchr/testify/assert"
)

const testUserName = "testUser"

func makeRestCall(t *testing.T, baseURL string, authHeader string) (int, string) {
	req, err := http.NewRequest("POST", baseURL+"checkToken", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making call: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(body))
	return resp.StatusCode, bodyStr
}

func checkTokenHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := netw.JwtGetTokenClaims(r.Context())
	if !ok {
		http.Error(w, "Error getting claims", http.StatusBadRequest)
		return
	}
	userName := claims["user"]
	if userName != testUserName {
		http.Error(w, "Unexpected user name "+userName, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "You just reached the checkToken handler")
}

func noToken(t *testing.T, baseURL string) {
	code, msg := makeRestCall(t, baseURL, "")

	assert.Equal(t, http.StatusUnauthorized, code)
	assert.Equal(t, "Authorization header missing", msg)
}

func noBearer(t *testing.T, baseURL string) {
	code, msg := makeRestCall(t, baseURL, "MiToken thereShouldHaveBeenBearer")

	assert.Equal(t, http.StatusUnauthorized, code)
	assert.Equal(t, "Invalid Authorization header format", msg)
}

func noBearer_OnePart(t *testing.T, baseURL string) {
	code, msg := makeRestCall(t, baseURL, "thereShouldHaveBeenBearer")

	assert.Equal(t, http.StatusUnauthorized, code)
	assert.Equal(t, "Invalid Authorization header format", msg)
}

func noBearer_SeveralParts(t *testing.T, baseURL string) {
	code, msg := makeRestCall(t, baseURL, "Bearer thereMustBeTokens whatAmIDoingHere")

	assert.Equal(t, http.StatusUnauthorized, code)
	assert.Equal(t, "Invalid Authorization header format", msg)
}

func tokenOk(t *testing.T, baseURL string) {
	token, err := jwt.GenerarToken(testUserName)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	code, msg := makeRestCall(t, baseURL, "Bearer "+token)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "You just reached the checkToken handler", msg)
}

func claimExtractionIsWorking(t *testing.T, baseURL string) {
	letsTrickTheTest := "ImNotTheUser"
	token, err := jwt.GenerarToken(letsTrickTheTest)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	code, msg := makeRestCall(t, baseURL, "Bearer "+token)

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, "Unexpected user name "+letsTrickTheTest, msg)
}

func TestMiddlewareJWT(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestMiddlewareJWT in short mode")
	}

	handlersFuncs := []libtest.HttpTestHandlerFunc{}
	handlers := []libtest.HttpTestHandler{
		{Path: "/checkToken", F: netw.JwtMiddleware(logger.GetNopLogger())(http.HandlerFunc(checkTokenHandler))},
	}

	testFunctions := []func(t *testing.T, baseURL string){
		noToken,
		noBearer,
		noBearer_OnePart,
		noBearer_SeveralParts,
		tokenOk,
		claimExtractionIsWorking,
	}

	libtest.RunSimpleServer(t, testFunctions, handlersFuncs, handlers)
}
