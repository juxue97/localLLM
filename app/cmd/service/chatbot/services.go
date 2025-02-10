package chatbot

import (
	"fmt"
	"net/http"

	"chatbot/config"
	"chatbot/utils"
)

func GenerateResponse(messages []map[string]string) (*http.Response, error) {
	// the request will be passed here, then if will be passed to curl
	// fmt.Println(JSONrequest.Query)

	payload := map[string]interface{}{
		"model":    config.Envs.LLMModel,
		"messages": messages,
		"stream":   true,
	}

	url := config.Envs.LLMIp + "/api/chat"
	// fmt.Println(payload)
	// curl: call the localLLM api and gather the response
	responseStreamBody, err := utils.CurlRequest(url, payload)
	if err != nil {
		return nil, fmt.Errorf("fail to curl, error : %w", err)
	}
	// defer responseStreamBody.Body.Close()
	return responseStreamBody, nil
}
