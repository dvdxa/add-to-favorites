package user_handler

import (
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	cases := []struct {
		name   string
		body   domain.User
		expErr error
	}{
		{
			name:   "invalid_length_of_user_credentials",
			body:   domain.User{Name: "zau", Password: "has"},
			expErr: ErrInvalidUserNameOrPassLength,
		},
		{
			name:   "invalid_characters_username",
			body:   domain.User{Name: "ddddd,", Password: "3fdfdcx"},
			expErr: ErrInvalidCharacters,
		},
		{
			name:   "invalid_characters_password",
			body:   domain.User{Name: "Khalid", Password: "//dkff,sddsads"},
			expErr: ErrInvalidCharacters,
		},
		{
			name:   "invalid_underscore_username",
			body:   domain.User{Name: "_khasan", Password: "fd23oxpd"},
			expErr: ErrInvalidUnderscore,
		},
		{
			name:   "invalid_underscore_pass",
			body:   domain.User{Name: "Khalid", Password: "_231311fdf"},
			expErr: ErrInvalidUnderscore,
		},
		{
			name:   "many_underscores",
			body:   domain.User{Name: "Val___id", Password: "2313fff"},
			expErr: ErrTooManyUnderscore,
		},
		{
			name:   "many_underscores_pass",
			body:   domain.User{Name: "Timbersaw", Password: "2131___fdsf"},
			expErr: ErrTooManyUnderscore,
		},
	}
	h := UserHandler{}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := h.ValidateRequest(tCase.body)
			if err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr.Error())
			}
		})
	}

}
