package main

import (
	"fmt"
	"testing"
)

func TestReplacingBadWords(t *testing.T) {
	tests := []struct{
		text	string
		result	string
	}{
		{
			text: "This is a kerfuffle opinion I need to share with the world",
			result: "This is a **** opinion I need to share with the world",
		},{
			text: "This is a kerfuff example",
			result: "This is a kerfuff example",
		},{
			text: "This is KerFuffLe! valid",
			result:  "This is KerFuffLe! valid",
		},{
			text: "This is KerFuffLe not valid",
			result:  "This is **** not valid",
		},{
			text: " this has to foRNax work, if    not I would not SHARBERT understand  ",
			result: " this has to **** work, if    not I would not **** understand  ",
		},
	}
	
	badWords := getBadWords()

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cleanText := replaceBadWords(test.text, badWordReplacement, badWords)
			if cleanText != test.result {
				t.Errorf("Expecting: '%s', obtained: '%s'", test.result, cleanText)
				return
			}
		})
	}

}
