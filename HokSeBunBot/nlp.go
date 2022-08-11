package main

import (
	"fmt"
	"log"
	"os"
	"path"

	hfapigo "github.com/TannerKvarfordt/hfapigo"
)

func init_nlp() {
	fmt.Println("Setting HF api:", CONFIG.HUGGINGFACE_TOKENs[0])
	log.Println("Setting HF api:", CONFIG.HUGGINGFACE_TOKENs[0])
	hfapigo.SetAPIKey(CONFIG.HUGGINGFACE_TOKENs[0])
}

func getSingleSummarization(filename string, input string) string {
	if _, err := os.Stat(path.Join(CONFIG.SUMMARIZATION_LOCATION, filename)); err == nil {
		bytes, err := os.ReadFile(path.Join(CONFIG.SUMMARIZATION_LOCATION, filename))
		if err != nil {
			log.Panicln(err)
		}
		return string(bytes)
	}

	sresps, err := hfapigo.SendSummarizationRequest(
		CONFIG.HUGGINGFACE_MODEL,
		&hfapigo.SummarizationRequest{
			Inputs:  []string{input},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		},
	)
	if err != nil {
		log.Panicln(err)
	}
	var content string = sresps[0].SummaryText

	// write summarization
	file, err := os.Create(path.Join(CONFIG.SUMMARIZATION_LOCATION, filename))
	if err != nil {
		log.Println(err)
	}
	file.WriteString(content)
	file.Close()
	log.Println("[HuggingFace] Get request for", filename, "content:", content)
	return content
}