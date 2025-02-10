package types

type Output struct {
	Response           string  `json:"response"`
	InputToken         int     `json:"inputToken"`
	OutputToken        int     `json:"outputToken"`
	LoadModelDuration  float64 `json:"loadModelDuration"`
	PromptEvalDuration float64 `json:"promptEvalDuration"`
	EvaluateDuration   float64 `json:"evaluateDuration"`
	TotalDuration      float64 `json:"totalDuration"`
}

type StreamResponse struct {
	Done    bool   `json:"done"`
	Message string `json:"message"`
	Data    Output `json:"data"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
