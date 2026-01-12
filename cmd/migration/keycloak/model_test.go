package keycloak

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap/pbkdf2"
	"github.com/zitadel/passwap/verifier"
)

func Test_user_getPassword(t *testing.T) {
	type fields struct {
		Credentials []credential
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "no password credential",
			fields: fields{
				Credentials: []credential{
					{Type: "foo"},
				},
			},
		},
		{
			name: "empty credential data",
			fields: fields{
				Credentials: []credential{
					{
						Type:       "password",
						SecretData: "{\"value\":\"ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=\",\"salt\":\"RaXjs4RiUKgJGkX6kp277w==\",\"additionalParameters\":{}}",
					},
				},
			},
		},
		{
			name: "empty secret data",
			fields: fields{
				Credentials: []credential{
					{
						Type:           "password",
						CredentialData: "{\"hashIterations\":27500,\"algorithm\":\"pbkdf2-sha256\",\"additionalParameters\":{}}",
					},
				},
			},
		},
		{
			name: "secret data error",
			fields: fields{
				Credentials: []credential{
					{
						Type:           "password",
						SecretData:     "~~~",
						CredentialData: "{\"hashIterations\":27500,\"algorithm\":\"pbkdf2-sha256\",\"additionalParameters\":{}}",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "credential data error",
			fields: fields{
				Credentials: []credential{
					{
						Type:           "password",
						SecretData:     "{\"value\":\"ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=\",\"salt\":\"RaXjs4RiUKgJGkX6kp277w==\",\"additionalParameters\":{}}",
						CredentialData: "~~~",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown algorithm",
			fields: fields{
				Credentials: []credential{
					{
						Type:           "password",
						SecretData:     "{\"value\":\"ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=\",\"salt\":\"RaXjs4RiUKgJGkX6kp277w==\",\"additionalParameters\":{}}",
						CredentialData: "{\"hashIterations\":27500,\"algorithm\":\"foobar\",\"additionalParameters\":{}}",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				Credentials: []credential{
					{
						Type:           "password",
						SecretData:     "{\"value\":\"ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=\",\"salt\":\"RaXjs4RiUKgJGkX6kp277w==\",\"additionalParameters\":{}}",
						CredentialData: "{\"hashIterations\":27500,\"algorithm\":\"pbkdf2-sha256\",\"additionalParameters\":{}}",
					},
				},
			},
			want: "$pbkdf2-sha256$27500$RaXjs4RiUKgJGkX6kp277w==$ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Credentials: tt.fields.Credentials,
			}
			got, err := u.getPassword()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			if tt.want != "" {
				// check if decoding works as expected.
				result, err := pbkdf2.NewVerifier(&pbkdf2.ValidationOpts{}).Verify(got, "x")
				require.NoError(t, err)
				assert.Equal(t, verifier.Fail, result)
			}
		})
	}
}
