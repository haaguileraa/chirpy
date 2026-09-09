package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		header		http.Header
		tokenToCompare	string
		expectingError	bool
		match		bool
	}{
		{
			header:		http.Header{	
				"Authorization" : {"Bearer tokenTest"},
			},
			tokenToCompare:	"tokenTest",
			expectingError:	false,
			match:		true,
		},{
			header:		http.Header{	
				"Authorization" : {"Bearer tokenTestWrong"},
			},
			tokenToCompare:	"tokenTest",
			expectingError:	false,
			match:		false,

		},{
			header:		http.Header{	
				"Authorization" : {"Bearer "},
			},
			tokenToCompare:	"tokentest",
			expectingError:	true,
			match:		false,
		},{
			header:		http.Header{	
				"Authorization" : {"Bearer"},
			},
			tokenToCompare:	"tokentest",
			expectingError:	true,
			match:		false,
		},{
			header:		http.Header{	
				"Authorization" : {"Bearer   "},
			},
			tokenToCompare:	"tokentest",
			expectingError:	true,
			match:		false,
		},{
			header:		http.Header{	
				"Authorization" : {"Bearer       tokenTest      "},
			},
			tokenToCompare:	"tokentest",
			expectingError:	true,
			match:		false,
		},{
			header:		http.Header{	
				"Authorization" : {"Bearer TOKENTEST           "},
			},
			tokenToCompare:	"TOKENTEST",
			expectingError:	false,
			match:		true,
		},{
			header:		http.Header{	
				"Content-Type" : {"application/json"},
			},
			tokenToCompare:	"tokentest",
			expectingError:	true,
			match:		false,
		},{
			header:		http.Header{	
				"Authorization" : {"BearerTokentest "},
			},
			tokenToCompare:	"Tokentest",
			expectingError:	true,
			match:		false,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			token, err := GetBearerToken(test.header)
			if err != nil { 
				if test.expectingError {
					return
				}
				t.Errorf("was not expecting error: %T", err)
				return
			}
			if token != test.tokenToCompare {
				if test.match {
					t.Errorf("was expecting a match '%s' != '%s'", token, test.tokenToCompare)
					return
				}
				return
			}
			if !test.match {
				t.Errorf("was not expecting a match '%s', '%s'", token, test.tokenToCompare)
			}
		})
	}
}
