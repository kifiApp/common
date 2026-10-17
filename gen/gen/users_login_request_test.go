package gen

import "testing"

func TestUserLoginRequestBiometricEnrollmentFields(t *testing.T) {
	request := &UserLoginRequest{
		EnableBiometric: true,
		PublicKeyJwk:    `{"kty":"EC","crv":"P-256","x":"x","y":"y"}`,
	}

	if !request.GetEnableBiometric() {
		t.Fatal("expected biometric enrollment to be enabled")
	}
	if request.GetPublicKeyJwk() == "" {
		t.Fatal("expected public JWK to be available")
	}
}
