package auth

import (
	"fmt"
	"testing"
)

func TestPasswords(t *testing.T) {
	tests := []struct {
		password	string
		password2Check	string
		match		bool
	}{
		{
			password: 	"bonnieIsADog",
			password2Check:	"thisIsNotRig",
			match:		false,
		},{
			password: 	"bonnieIsADog",
			password2Check:	"bonnieIsADog",
			match:		true,
		},{
			password: 	"_bonnieIsADog123",
			password2Check:	"bonnieIsADog",
			match:		false,
		},{
			password: 	"_bonnieIsADog123",
			password2Check:	"_bonnieIsADog123",
			match:		true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			hash, err := HashPassword(test.password)
			if err != nil {
				t.Errorf("Not expecting an error on Hash, but: %T", err)
				return
			}
			match, err := CheckPasswordHash(test.password2Check, hash)
			if err != nil {
				t.Errorf("Not expecting an error checking password, but: %T", err)
				return
			}
			if match != test.match {
				t.Errorf("Error checking password, expecting match: '%v', got '%v'", test.match, match)
				return
			}

		})
	}

}
