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
	err = subt.WriteToFile()
	if err != nil {
		t.Fatal("subt.writeToFile(): Error writing files")
	}

	// read from JSON
	var subt2 SubtitleSRT
	err = subt2.ReadFromFile()
	if err != nil {
		t.Fatal("subt.readFromFile(): Error reading files")
	}

	if !(subt.isEqual(subt2)) {
		t.Fatal("Subt IO: Read object is different from original object")
	}
}

func TestMoveWordToPrev(t *testing.T) {
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

	// Move word from line 1 to 0
	subt.MoveWordFromLineToPrev(1)
	want := "Hola a todos ¿cómo"
	if subt.translatedLine[0] != want {
		t.Fatalf("MoveWordFromLineToPrev(1): Line 0: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
	want = "estáis?"
	if subt.translatedLine[1] != want {
		t.Fatalf("MoveWordFromLineToPrev(1): Line 1: Want '%s' have '%s'", want, subt.translatedLine[1])
	}
}

func TestMoveWordToPrevSingleWord(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from line 13 to 12
	subt.MoveWordFromLineToPrev(13)
	want = ""
	if subt.translatedLine[13] != want {
		t.Fatalf("MoveWordFromLineToPrev(13): Line 13: Want '%s' have '%s'", want, subt.translatedLine[13])
	}
	want = "otra línea exacta exacta"
	if subt.translatedLine[12] != want {
		t.Fatalf("MoveWordFromLineToPrev(13): Line 12: Want '%s' have '%s'", want, subt.translatedLine[12])
	}
}

func TestMoveWordFrom0toNext(t *testing.T) {
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

	// Move word from line 0 to 1
	subt.MoveWordFromLineToNext(0)
	want := "Hola a"
	if subt.translatedLine[0] != want {
		t.Fatalf("MoveWordFromLineToNext(0): Line 0: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
	want = "todos ¿cómo estáis?"
	if subt.translatedLine[1] != want {
		t.Fatalf("MoveWordFromLineToNext(0): Line 1: Want '%s' have '%s'", want, subt.translatedLine[1])
	}
}

func TestMoveWordToNextSingleWord(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from line 9 to 10
	subt.MoveWordFromLineToNext(9)
	want = ""
	if subt.translatedLine[9] != want {
		t.Fatalf("MoveWordFromLineToPrev(9): Line 9: Want '%s' have '%s'", want, subt.translatedLine[9])
	}
	want = "hoy?"
	if subt.translatedLine[10] != want {
		t.Fatalf("MoveWordFromLineToPrev(10): Line 10: Want '%s' have '%s'", want, subt.translatedLine[10])
	}
}

func TestMoveWordToPreviousLine0(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from line 0 to previous, check that nothing happened
	subt.MoveWordFromLineToPrev(0)
	want = "Hola a todos"
	if subt.translatedLine[0] != want {
		t.Fatalf("MoveWordFromLineToPrev(0): Line 0: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
}

func TestMoveWordToNextLine15(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from last line to next, check that nothing happened
	subt.MoveWordFromLineToNext(15)
	want = "y gracias"
	if subt.translatedLine[15] != want {
		t.Fatalf("MoveWordFromLineToNext(15): Line 15: Want '%s' have '%s'", want, subt.translatedLine[15])
	}
}

func TestMoveWordLineToNextLine1(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from line 1 to next (last of lineset), check that nothing happened
	subt.MoveWordFromLineToNext(1)
	want = "¿cómo estáis?"
	if subt.translatedLine[1] != want {
		t.Fatalf("MoveWordFromLineToNext(1): Line 1: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
}

func TestMoveWordLineToPrevLine14(t *testing.T) {
	fileNameSrt := "../datasubt/en.srt"
	fileNameTxt := "../datasubt/es.txt"
	var subt SubtitleSRT
	var want string

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.ImportOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.ImportTranslatedText(PrepareString(string(content)))

	// Move word from line 11 to previous (first of lineset), check that nothing happened
	subt.MoveWordFromLineToPrev(11)
	want = "Una línea exacta"
	if subt.translatedLine[11] != want {
		t.Fatalf("MoveWordFromLineToPrev(8): Line 8: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
}

func TestIsNotLoaded(t *testing.T) {
	var subt SubtitleSRT

	isLoaded := subt.IsLoaded()

	if isLoaded {
		t.Fatal("IsLoaded(): Not isLoaded test, Want false have true")
	}
}

func TestIsLoaded(t *testing.T) {
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

	isLoaded := subt.IsLoaded()

	if !isLoaded {
		t.Fatal("IsLoaded():  isLoaded test, Want true have false")
	}
}
