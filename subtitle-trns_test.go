package subtitle

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestTranslateOrigText(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// Translate the original text
	subt.Translate("es", "laboratorio-cloudmas-efraim", "nmt")

	// Verify the translation is consistent
	if !subt.IsTranslationConsistent() {
		t.Fatalf("Translation failed: Translation is not consistent")
	}
	// Verify the content
	want := []string{
		"Hola a todos",
		"¿Cómo están?",
		"[Música]",
		"[Más música]",
		"",
		"[Música]",
		"",
		"Hola a todos",
		"¿Cómo están",
		"hoy?",
		"",
		"Una línea exacta",
		"otra línea exacta",
		"exacta",
		"Adiós a todos",
		"y gracias",
	}
	have := subt.GetTranslatedLines()
	isEqual := true
	if len(want) != len(have) {
		isEqual = false
	} else {
		for i, v := range have {
			if v != want[i] {
				isEqual = false
				break
			}
		}
	}
	if !isEqual {
		t.Fatalf("Translation failed: want %q have %q.", want, subt.translatedLine)
	}
}
