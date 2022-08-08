package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	MassResponse struct {
		Codes []int  `validate:"in:201,403,502"`
		Body  string `json:"omitempty"`
	}

	SomeRoles struct {
		Roles []string `validate:"in:admin,user"`
	}

	BadRule struct {
		Version string `validate:"bad:rule"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "value_is_not_a_struct",
			expectedErr: errors.New("value is not a struct"),
		},
		{
			in: User{
				ID:    "1",
				Name:  "noname",
				Age:   10,
				Email: "it`snot.email",
				Phones: []string{
					"34634534765",
					"09898",
				},
				Role: "stufffff",
			},
			expectedErr: errors.New("error in the field: ID - the length of the line must be exactly 36\n" +
				"error in the field: Age - the number cannot be less than 18\n" +
				"error in the field: Email - the string must match the regular expression\n" +
				"error in the field: Role - stufffff must be part of the set admin,stuff\n" +
				"error in the field: Phones - the length of the line must be exactly 11"),
		},
		{
			in: User{
				ID:    "123456789_123456789_123456789_123456789_1234567890",
				Name:  "noname",
				Age:   100,
				Email: "it's not email",
				Phones: []string{
					"123456789_1234567890",
				},
				Role: "guest",
			},
			expectedErr: errors.New("error in the field: ID - the length of the line must be exactly 36\n" +
				"error in the field: Age - the number cannot be greater than 50\n" +
				"error in the field: Email - the string must match the regular expression\n" +
				"error in the field: Role - guest must be part of the set admin,stuff\n" +
				"error in the field: Phones - the length of the line must be exactly 11"),
		},
		{
			in:          App{Version: "123456789"},
			expectedErr: errors.New("error in the field: Version - the length of the line must be exactly 5"),
		},
		{
			in:          App{Version: ""},
			expectedErr: errors.New("error in the field: Version - the length of the line must be exactly 5"),
		},
		{
			in:          Response{Code: 451},
			expectedErr: errors.New("error in the field: Code - 451 must be part of the set 200,404,500"),
		},
		{
			in:          MassResponse{Codes: []int{451, 404, 500}},
			expectedErr: errors.New("error in the field: Codes - 451 must be part of the set 201,403,502"),
		},
		{
			in:          SomeRoles{Roles: []string{"supermanager", "manager"}},
			expectedErr: errors.New("error in the field: Roles - supermanager must be part of the set admin,user"),
		},

		{
			in:          BadRule{Version: "1.16"},
			expectedErr: errors.New("rule 'Bad' not exist"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

func TestPassValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "123456789_123456789_123456789_123456",
				Name:  "Evgenii",
				Age:   30,
				Email: "public@softroot.ru",
				Role:  "admin",
				Phones: []string{
					"89100000000",
				},
				meta: nil,
			},
		},
		{
			in: Token{},
		},
		{
			in: Response{Code: 500},
		},
		{
			in: Response{Code: 404},
		},
		{
			in: Response{Code: 200},
		},
		{
			in: MassResponse{Codes: []int{201, 403}},
		},
		{
			in: MassResponse{Codes: []int{201, 502, 403}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			require.Nil(t, err)
		})
	}
}
