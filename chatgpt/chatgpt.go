package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xmodules/chatgpt/assets"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	/*
	   MaxTokens   int       `json:"max_tokens"`
	   Top_p       int       `json:"top_p"`
	   Stream      bool      `json:"stream"`
	   N           int       `json:"n"`
	   Stop        string    `json:"stop"`
	*/
}

func Ask(ds applications.Datasource, data string) (string, error) {

	Msg := []Message{{"user", data}}

	cfg := ds.GetConfig().GetConfig(assets.MODULEID)
	if cfg == nil {
		return "", fmt.Errorf("Module %s not configured into the datasource", assets.MODULEID)
	}

	m, _ := cfg.GetString("model")
	ep, _ := cfg.GetString("endpoint")
	ak, _ := cfg.GetString("apikey")

	reqBody := CompletionRequest{
		Model:       m,
		Messages:    Msg,
		Temperature: 0.7,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ak))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Translate(ds applications.Datasource, data string, lang string) (string, error) {

	return Ask(ds, "Traduce las líneas siguientes al "+lang+", para una página web de recetas de cocina, guardando el mismo formato y sin quitar los números ni los | y sin poner un prompt en la respuesta:\n"+data)
}

func TranslatePrompt(ds applications.Datasource, prompt string, data string) (string, error) {

	return Ask(ds, prompt+"\n"+data)
}
