package i18y

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {
	err := Init()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		acceptLanguage string
		key            string
		out            string
	}{
		{
			name:           "japanese",
			acceptLanguage: language.Japanese.String(),
			key:            "cost",
			out:            "料金",
		},
		{
			name:           "english",
			acceptLanguage: language.English.String(),
			key:            "cost",
			out:            "Cost",
		},
	}

	for _, tt := range tests {
		out := Translate(tt.acceptLanguage, tt.key)
		assert.Equal(t, tt.out, out)
	}
}
