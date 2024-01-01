package i18y

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

const (
	LanguagesDir = "languages"
)

var (
	languages = []language.Tag{
		language.Japanese,
		language.English,
	}
	//go:embed languages/*.yaml
	configByte embed.FS
)

type Messages map[string]map[string]string

func Init() error {
	files, err := configByte.ReadDir(LanguagesDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := fs.ReadFile(configByte, fmt.Sprintf("%s/%s", LanguagesDir, file.Name()))
		if err != nil {
			return err
		}

		var m map[string]string
		err = yaml.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		// Make language tag from file name
		l := language.Make(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		for k, v := range m {
			err = message.SetString(l, k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Translate(acceptLanguage string, msg string, args ...interface{}) string {
	t, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return msg
	}

	matcher := language.NewMatcher(languages)
	tag, _, _ := matcher.Match(t...)
	p := message.NewPrinter(tag)
	return p.Sprintf(msg, args...)
}
