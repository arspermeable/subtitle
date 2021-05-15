package subtitle

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestWriteReadFile(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// write to JSON
	err = subt.writeToFile()
	if err != nil {
		t.Fatal("subt.writeToFile(): Error writing files")
	}

	// read from JSON
	var subt2 SubtitleSRT
	err = subt2.readFromFile()
	if err != nil {
		t.Fatal("subt.readFromFile(): Error reading files")
	}

	if !(subt.isEqual(subt2)) {
		t.Fatal("Subt IO: Read object is different from original object")
	}

	t.Log("The original and read objects are equal!")
}
