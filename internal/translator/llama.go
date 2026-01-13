package translator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fdddf/xcstrings-translator/internal/model"
	"github.com/gofiber/fiber/v2/log"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// globalLlamaMutex serializes all llama operations to prevent race conditions
var globalLlamaMutex sync.Mutex

// LlamaTranslator implements the TranslationProvider interface for local Llama models
type LlamaTranslator struct {
	ModelPath        string
	GrammarPath      string
	Threads          int
	Seed             int
	Tokens           int
	TopK             int
	Tfs              float64
	TopP             float64
	TypicalP         float64
	RepeatLastN      int
	RepeatPenalty    float64
	FrequencyPenalty float64
	PresencePenalty  float64
	Temperature      float64
	Verbose          bool
	Model            llama.Model
	ContextParams    llama.ContextParams
	SamplerParams    llama.SamplerParams
	Samplers         []llama.SamplerType
	mutex            sync.Mutex
}

// LlamaOptions represents configuration options for the Llama translator
type LlamaOptions struct {
	ModelPath        string
	GrammarPath      string
	Threads          int
	Seed             int
	Tokens           int
	TopK             int
	Tfs              float64
	TopP             float64
	TypicalP         float64
	RepeatLastN      int
	RepeatPenalty    float64
	FrequencyPenalty float64
	PresencePenalty  float64
	Temperature      float64
	Verbose          bool
}

// NewLlamaTranslator creates a new Llama Translator instance
func NewLlamaTranslator(options LlamaOptions) (*LlamaTranslator, error) {
	// Note: The llama library should already be loaded and initialized by the calling code
	// as shown in cmd/llama.go. We'll just check if we can load the model.

	if !options.Verbose {
		llama.LogSet(llama.LogSilent())
	}

	modelParams := llama.ModelDefaultParams()

	model, err := llama.ModelLoadFromFile(options.ModelPath, modelParams)
	if err != nil {
		return nil, fmt.Errorf("failed to load llama model from %s: %v", options.ModelPath, err)
	}
	if model == 0 {
		return nil, fmt.Errorf("failed to load llama model from %s", options.ModelPath)
	}

	// Create context parameters
	ctxParams := llama.ContextDefaultParams()
	if options.Threads > 0 {
		ctxParams.NThreads = int32(options.Threads)
		ctxParams.NThreadsBatch = int32(options.Threads)
	}

	// Initialize sampler with the specified parameters
	samplers := []llama.SamplerType{
		llama.SamplerTypeTopK,
		llama.SamplerTypeTopP,
		llama.SamplerTypeMinP,
		llama.SamplerTypeTemperature,
		llama.SamplerTypePenalties,
	}

	// Add sampler types based on configuration
	// if options.TypicalP > 0 && options.TypicalP < 1.0 {
	// 	samplers = append(samplers, llama.SamplerTypeTypicalP)
	// }
	// if options.RepeatLastN > 0 && (options.RepeatPenalty != 1.0 || options.FrequencyPenalty != 0.0 || options.PresencePenalty != 0.0) {
	// 	samplers = append(samplers, llama.SamplerTypePenalties)
	// }

	samplerParams := llama.DefaultSamplerParams()
	samplerParams.Temp = float32(options.Temperature)
	samplerParams.TopK = int32(options.TopK)
	samplerParams.TopP = float32(options.TopP)
	samplerParams.TypP = float32(options.TypicalP)
	samplerParams.PenaltyLastN = int32(options.RepeatLastN)
	samplerParams.PenaltyRepeat = float32(options.RepeatPenalty)
	samplerParams.PenaltyFreq = float32(options.FrequencyPenalty)
	samplerParams.PenaltyPresent = float32(options.PresencePenalty)

	// Create the translator instance
	translator := &LlamaTranslator{
		ModelPath:        options.ModelPath,
		GrammarPath:      options.GrammarPath,
		Threads:          options.Threads,
		Seed:             options.Seed,
		Tokens:           options.Tokens,
		TopK:             options.TopK,
		Tfs:              options.Tfs,
		TopP:             options.TopP,
		TypicalP:         options.TypicalP,
		RepeatLastN:      options.RepeatLastN,
		RepeatPenalty:    options.RepeatPenalty,
		FrequencyPenalty: options.FrequencyPenalty,
		PresencePenalty:  options.PresencePenalty,
		Temperature:      options.Temperature,
		Verbose:          options.Verbose,
		Model:            model,
		ContextParams:    ctxParams,
		SamplerParams:    *samplerParams,
		Samplers:         samplers,
	}

	return translator, nil
}

// Translate translates a string using local Llama model
func (l *LlamaTranslator) Translate(ctx context.Context, req model.TranslationRequest) (model.TranslationResponse, error) {
	// Prepare translation prompt
	prompt := createTranslationPrompt(req)

	messages := make([]llama.ChatMessage, 0)
	messages = append(messages, llama.NewChatMessage("user", prompt))

	template := llama.ModelChatTemplate(l.Model, "")
	if template == "" {
		template = "chatml"
	}
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(template, messages, true, buf)
	result := string(buf[:len])

	// Generate translation using the Llama model (with mutex safety in generateText)
	response, err := l.generateText(ctx, result)
	if err != nil {
		return model.TranslationResponse{
			Key:            req.Key,
			TargetLanguage: req.TargetLanguage,
			Error:          fmt.Errorf("llama generation failed: %v", err),
		}, nil
	}

	return model.TranslationResponse{
		Key:            req.Key,
		TargetLanguage: req.TargetLanguage,
		TranslatedText: cleanLlamaResponse(response),
	}, nil
}

// generateText generates text using the Llama model based on the provided prompt
func (l *LlamaTranslator) generateText(ctx context.Context, prompt string) (string, error) {
	// Use global mutex to ensure thread safety across all llama operations
	globalLlamaMutex.Lock()
	defer globalLlamaMutex.Unlock()

	// Create a new context for this generation
	llamaContext, err := llama.InitFromModel(l.Model, l.ContextParams)
	if err != nil {
		return "", fmt.Errorf("failed to initialize llama context: %v", err)
	}
	defer llama.Free(llamaContext)

	// Create sampler params for this request
	samplerParams := l.SamplerParams
	// Initialize sampler for this request
	sampler := llama.NewSampler(l.Model, l.Samplers, &samplerParams)
	defer llama.SamplerFree(sampler)

	// Get the vocabulary from the model
	vocab := llama.ModelGetVocab(l.Model)

	// Tokenize the prompt
	tokens := llama.Tokenize(vocab, prompt, true, true)

	// Check if we have tokens to process
	if len(tokens) == 0 {
		return "", fmt.Errorf("no tokens generated from prompt")
	}

	// Process the prompt tokens
	batch := llama.BatchGetOne(tokens)

	// Encode the tokens in the batch
	if llama.ModelHasEncoder(l.Model) {
		_, err := llama.Encode(llamaContext, batch)
		if err != nil {
			return "", fmt.Errorf("failed to encode batch: %v", err)
		}

		start := llama.ModelDecoderStartToken(l.Model)
		if start == llama.TokenNull {
			start = llama.VocabBOS(vocab)
		}

		batch = llama.BatchGetOne([]llama.Token{start})
	} else {
		_, err := llama.Decode(llamaContext, batch)
		if err != nil {
			return "", fmt.Errorf("failed to decode batch: %v", err)
		}
	}

	// Generate tokens
	maxTokens := l.Tokens
	if maxTokens <= 0 {
		maxTokens = 512 // default
	}

	var result strings.Builder
	for pos := int32(0); pos < int32(maxTokens); pos += batch.NTokens {
		// Decode the current context
		_, err := llama.Decode(llamaContext, batch)
		if err != nil {
			return "", fmt.Errorf("failed to decode: %v", err)
		}

		// Sample next token
		nextToken := llama.SamplerSample(sampler, llamaContext, -1)

		// Check if it's an end-of-sequence token
		if llama.VocabIsEOG(vocab, nextToken) {
			break
		}

		// Convert token to text piece
		buf := make([]byte, 256)
		l := llama.TokenToPiece(vocab, nextToken, buf, 0, false)
		nextPiece := string(buf[:l])

		// Add to result
		result.WriteString(nextPiece)

		// Create new batch with the next token
		batch = llama.BatchGetOne([]llama.Token{nextToken})
	}

	return result.String(), nil
}

// createTranslationPrompt creates a prompt for translating text
func createTranslationPrompt(req model.TranslationRequest) string {
	// Create a descriptive prompt for translation task
	prompt := fmt.Sprintf(
		"Translate the following segment into %s, without additional explanation.\n\n%s",
		// languageDisplayName(req.SourceLanguage),
		languageDisplayName(req.TargetLanguage),
		req.Text,
	)
	log.Infof("prompt: %s", prompt)

	return prompt
}

func languageDisplayName(code string) string {
	if code == "" {
		return code
	}

	nameByCode := map[string]string{
		"ar":      "Arabic",
		"cs":      "Czech",
		"da":      "Danish",
		"de":      "German",
		"en":      "English",
		"es":      "Spanish",
		"fi":      "Finnish",
		"fr":      "French",
		"he":      "Hebrew",
		"hi":      "Hindi",
		"id":      "Indonesian",
		"it":      "Italian",
		"ja":      "Japanese",
		"ko":      "Korean",
		"nl":      "Dutch",
		"no":      "Norwegian",
		"pl":      "Polish",
		"pt":      "Portuguese",
		"ru":      "Russian",
		"sv":      "Swedish",
		"th":      "Thai",
		"tr":      "Turkish",
		"vi":      "Vietnamese",
		"zh-hans": "Chinese (Simplified)",
		"zh-hant": "Chinese (Traditional)",
	}

	normalized := strings.ToLower(strings.TrimSpace(code))
	if name, ok := nameByCode[normalized]; ok {
		return name
	}
	return code
}

func cleanLlamaResponse(response string) string {
	trimmed := strings.TrimSpace(response)
	if trimmed == "" {
		return trimmed
	}

	upper := strings.ToUpper(trimmed)
	markers := []string{"TRANSLATION:", "TEXT:"}
	for _, marker := range markers {
		if idx := strings.Index(upper, marker); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[idx+len(marker):])
			upper = strings.ToUpper(trimmed)
		}
	}

	if strings.Contains(trimmed, "\n") {
		lines := strings.Split(trimmed, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			trimmed = line
			break
		}
	}

	return strings.TrimSpace(trimmed)
}

// Close closes the Llama translator resources
func (l *LlamaTranslator) Close() {
	// The mutex doesn't need to be locked here since this should only be called once
	// when the translator is being closed
	if l.Model != 0 {
		llama.ModelFree(l.Model)
	}
	// llama.Close() will be called externally after all translators are closed
}
