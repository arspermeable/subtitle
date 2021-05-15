package subtitle

import (
	"encoding/json"
	"io/ioutil"
)

// Read from/Write to JSON files

// writeToFile writes a SubtitleSRT struct into JSON files
// Each field will be saved in a separate file
//    SubtitleBlock -> subtitleblock.json
//    LineSet -> lineset.json
//    originalLine -> originalline.json
//    translatedLine -> translatedline.json
//    translatedSet -> translatedset.json
//    translatedText -> translatedtext.json
func (this *SubtitleSRT) WriteToFile() error {
	var data []byte
	var err error

	// Save SubtitleBlock
	data, err = json.MarshalIndent(this.subtitleBlock, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("subtitleblock.json", data, 0644)
	if err != nil {
		return err
	}

	// Save LineSet
	data, err = json.MarshalIndent(this.lineSet, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("lineset.json", data, 0644)
	if err != nil {
		return err
	}

	// Save originalLine
	data, err = json.MarshalIndent(this.originalLine, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("originalline.json", data, 0644)
	if err != nil {
		return err
	}

	// Save translatedLine
	data, err = json.MarshalIndent(this.translatedLine, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("translatedline.json", data, 0644)
	if err != nil {
		return err
	}

	// Save translatedSet
	data, err = json.MarshalIndent(this.translatedSet, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("translatedset.json", data, 0644)
	if err != nil {
		return err
	}

	// Save translatedText
	data, err = json.MarshalIndent(this.translatedText, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("translatedtext.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// writeToFile writes a SubtitleSRT struct into JSON files
// Each field will be saved in a separate file
//    SubtitleBlock <- subtitleblock.json
//    LineSet <- lineset.json
//    originalLine <- originalline.json
//    translatedLine <- translatedline.json
//    translatedSet <- translatedset.json
//    translatedText <- translatedtext.json
func (this *SubtitleSRT) ReadFromFile() error {
	var data []byte
	var err error

	// Read SubtitleBlock
	data, err = ioutil.ReadFile("subtitleblock.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.subtitleBlock)
	if err != nil {
		return err
	}

	// Read LineSet
	data, err = ioutil.ReadFile("lineset.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.lineSet)
	if err != nil {
		return err
	}

	// Read originalLine
	data, err = ioutil.ReadFile("originalline.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.originalLine)
	if err != nil {
		return err
	}

	// Read translatedLine
	data, err = ioutil.ReadFile("translatedline.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.translatedLine)
	if err != nil {
		return err
	}

	// Read translatedSet
	data, err = ioutil.ReadFile("translatedset.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.translatedSet)
	if err != nil {
		return err
	}

	// Read translatedText
	data, err = ioutil.ReadFile("translatedtext.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &this.translatedText)
	if err != nil {
		return err
	}

	return nil
}
