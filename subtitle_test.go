package subtitle

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestSetAllData1(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Verify the "controlled environment" worked well
	want := []string{
		"Hola a todos",
		"¿cómo estáis?",
		"[Música]",
		"[Más música]",
		"",
		"[Música]",
		"",
		"Hola a \"todos\",",
		"¿cómo estáis",
		"hoy?",
		"",
		"Una línea exacta",
		"otra línea exacta",
		"exacta",
		"Adios a todos",
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
		t.Fatalf("Import SRT/TRT failed: want %q, have %q", want, subt.translatedLine)
	}

}

func TestSetAllData1b(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1b.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Verify the "controlled environment" worked well
	want := []string{
		"Hola a todos",
		"¿cómo estáis?",
		"[Música]",
		"[Más música]",
		"",
		"[Música]",
		"",
		"Hola a \"todos\",",
		"¿cómo estáis",
		"hoy?",
		"",
		"Una línea exacta",
		"otra línea exacta",
		"exacta",
		"Adios a todos",
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
		t.Fatalf("Import SRT/TRT failed: want %q, have %q", want, subt.translatedLine)
	}

}
func TestWriteReadFile(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

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

	if !(subt.IsEqual(subt2)) {
		t.Fatal("Subt IO: Read object is different from original object")
	}
}

func TestMoveWordFromLineSetToPrev_ManyWords(t *testing.T) {
	var fileNameSrt string
	var fileNameTxt string
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from line 3 to prev (lineSet[3] has 6 words)
	subt.MoveWordsFromLineSetToPrev(3, 100)
	want := []string{
		"Hola a todos",
		"¿cómo estáis?",
		"[Música]",
		"[Más música]",
		"",
		"[Música] Hola a \"todos\", ¿cómo estáis hoy?",
		"",
		"",
		"",
		"",
		"",
		"Una línea exacta",
		"otra línea exacta",
		"exacta",
		"Adios a todos",
		"y gracias",
	}

	for i, str := range want {
		if str != subt.translatedLine[i] {
			t.Fatalf("MoveWordsFromLineToPrev(3,100): Line %d: Want '%s' have '%s'", i, str, subt.translatedLine[i])
		}
	}
	if !subt.IsTranslationConsistent() {
		t.Fatalf("MoveWordsFromLineSetToPrev(3,100): Translation is not consistent after movement")

	}
}

func TestMoveWordFromLineSetToNext_ManyWords(t *testing.T) {
	var fileNameSrt string
	var fileNameTxt string
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from line 3 to prev (lineSet[3] has 6 words)
	subt.MoveWordsFromLineSetToNext(3, 100)
	want := []string{
		"Hola a todos",
		"¿cómo estáis?",
		"[Música]",
		"[Más música]",
		"",
		"[Música]",
		"",
		"",
		"",
		"",
		"",
		"Hola a \"todos\", ¿cómo estáis",
		"hoy? Una línea exacta otra línea",
		"exacta exacta",
		"Adios a todos",
		"y gracias",
	}

	for i, str := range want {
		if str != subt.translatedLine[i] {
			t.Fatalf("MoveWordsFromLineToNext(3,100): Line %d: Want '%s' have '%s'", i, str, subt.translatedLine[i])
		}
	}
	if !subt.IsTranslationConsistent() {
		t.Fatalf("MoveWordsFromLineSetToNext(3,100): Translation is not consistent after movement")

	}
}

func TestMoveWordToPrev(t *testing.T) {
	var fileNameSrt string
	var fileNameTxt string
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

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

func TestMoveWordToPrevFirstLine(t *testing.T) {
	var fileNameSrt string
	var fileNameTxt string
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from line 1 to 0
	subt.MoveWordFromLineToPrev(14)
	want := []string{"Una línea exacta", "otra línea exacta exacta", "Adios", "a todos y", "gracias"}
	for i, str := range want {
		if str != subt.translatedLine[11+i] {
			t.Fatalf("MoveWordFromLineToPrev(14): Line %d: Want '%s' have '%s'", 11+i, str, subt.translatedLine[11+i])
		}
	}
	if !subt.IsTranslationConsistent() {
		t.Fatalf("MoveWordFromLineToPrev(14): Translation is not consistent after movement")

	}
}

func TestMoveWordToNextLastLine(t *testing.T) {
	var fileNameSrt string
	var fileNameTxt string
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from line 1 to 0
	subt.MoveWordFromLineToNext(13)
	want := []string{"Una línea exacta", "otra línea", "exacta", "exacta Adios a", "todos y gracias"}
	for i, str := range want {
		if str != subt.translatedLine[11+i] {
			t.Fatalf("MoveWordFromLineToNext(13): Line %d: Want '%s' have '%s'", 11+i, str, subt.translatedLine[11+i])
		}
	}
	if !subt.IsTranslationConsistent() {
		t.Fatalf("MoveWordFromLineToNext(13): Translation is not consistent after movement")

	}
}

func TestMoveWordToPrevSingleWord(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT
	var want string

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

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
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

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
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT
	var want string

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

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
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT
	var want string

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from line 0 to previous, check that nothing happened
	subt.MoveWordFromLineToPrev(0)
	want = "Hola a todos"
	if subt.translatedLine[0] != want {
		t.Fatalf("MoveWordFromLineToPrev(0): Line 0: Want '%s' have '%s'", want, subt.translatedLine[0])
	}
}

func TestMoveWordToNextLine15(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT
	var want string

	// This test function uses simple test data
	fileNameSrt = "../datasubt/test1.srt"
	fileNameTxt = "../datasubt/test1.txt"

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	// Move word from last line to next, check that nothing happened
	want = "y gracias"
	if subt.translatedLine[15] != want {
		t.Fatalf("MoveWordFromLineToNext(15): Line 15: Want '%s' have '%s'", want, subt.translatedLine[15])
	}
	subt.MoveWordFromLineToNext(15)
	want = "y gracias"
	if subt.translatedLine[15] != want {
		t.Fatalf("MoveWordFromLineToNext(15): Line 15: Want '%s' have '%s'", want, subt.translatedLine[15])
	}
}

func TestDeleteAllData(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	subt.DeleteSubtitleSrt()
	if subt.IsLoaded() {
		t.Fatal("DeleteSubtitleSrt(): Data hasn't been deleted, IsLoaded()=true")
	}
}

func TestIsNotLoaded(t *testing.T) {
	var subt SubtitleSRT

	isLoaded := subt.IsLoaded()
	if isLoaded {
		t.Fatal("IsLoaded(): Want false have true")
	}
}

func TestIsLoaded(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	isLoaded := subt.IsLoaded()

	if !isLoaded {
		t.Fatal("IsLoaded():  isLoaded test, Want true have false")
	}
}

func TestTranslationConsistency(t *testing.T) {
	var fileNameSrt string = "../datasubt/test1.srt"
	var fileNameTxt string = "../datasubt/test1.txt"
	var subt SubtitleSRT

	// Import the subtitle file
	data, err := ioutil.ReadFile(fileNameSrt)
	check(err)
	reader := strings.NewReader(string(data))
	subt.SetOriginalSrt(reader)

	// import the translated text
	content, err := ioutil.ReadFile(fileNameTxt)
	check(err)
	subt.SetTranslatedText(string(content))

	if !subt.IsTranslationConsistent() {
		t.Fatal("IsTranslationConsistent(): Returns false, wanted true")
	}
}
