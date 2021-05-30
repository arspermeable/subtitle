package subtitle

import (
	"context"
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// Translate() translates the original text in originalLine
// into the requested language.
// Then, it stores the translatedText and splits line sets and translatedLine
func (this *SubtitleSRT) Translate(targetLang string, projectID string, model string) int {
	// Verify that data is already loaded
	if !this.IsLoadedSRT() {
		return 0
	}

	// Get a context
	ctx := context.Background()
	// Create a new translation client
	// No credentials are provided, this must be executed with
	// $GOOGLE_APPLICATION_CREDENTIALS correctly defined, as in /etc/environment
	client, err := translate.NewTranslationClient(ctx)
	check(err)
	defer client.Close()

	// This is for Batch request...
	/*
		req := &translatepb.BatchTranslateTextRequest{
			Parent:              fmt.Sprintf("projects/%s/locations/%s", projectID, location),
			SourceLanguageCode:  sourceLang,
			TargetLanguageCodes: []string{destLang},
			Models:              map[string]string{},
			InputConfigs:        []*translatepb.InputConfig{},
			OutputConfig:        &translatepb.OutputConfig{},
			Glossaries:          map[string]*translatepb.TranslateTextGlossaryConfig{},
			Labels:              map[string]string{},
		}

		op, err := client.BatchTranslateText(ctx, req)
		if err != nil {
			return "", err
		}

		resp, err := op.Wait(ctx)
		if err != nil {
			return "", err
		}
	*/

	txt, _ := this.GetOriginalText()

	req := &translatepb.TranslateTextRequest{
		Contents:           []string{txt},
		MimeType:           "text/plain",
		SourceLanguageCode: "en",
		TargetLanguageCode: targetLang,
		Parent:             fmt.Sprintf("projects/%s", projectID),
		Model:              fmt.Sprintf("projects/%s/locations/global/models/general/%s", projectID, model),
	}

	resp, err := client.TranslateText(ctx, req)
	check(err)

	// Store the translatedText
	this.SetTranslatedText(resp.GetTranslations()[0].GetTranslatedText())

	return len([]rune(this.translatedText))

}
