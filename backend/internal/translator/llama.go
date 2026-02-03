package translator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fdddf/opentrans/internal/model"
	hunyuanLang "github.com/fdddf/opentrans/pkg/hunyuan"
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
	MinP             float64
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
	ctxParams.NCtx = uint32(options.Tokens)
	ctxParams.NBatch = uint32(options.Tokens)
	ctxParams.NUbatch = uint32(options.Tokens)

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
	samplerParams.MinP = float32(options.MinP)

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
	buf := make([]byte, 4096)
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
	log.Infof("translate result: %s", response)

	return model.TranslationResponse{
		Key:            req.Key,
		TargetLanguage: req.TargetLanguage,
		TranslatedText: cleanLlamaResponse(response),
	}, nil
}

// TranslateStream translates a string using local Llama model with streaming support
func (l *LlamaTranslator) TranslateStream(ctx context.Context, req model.TranslationRequest, callback func(chunk string) error) (model.TranslationResponse, error) {
	// Prepare translation prompt
	prompt := createTranslationPrompt(req)

	messages := make([]llama.ChatMessage, 0)
	messages = append(messages, llama.NewChatMessage("user", prompt))

	template := llama.ModelChatTemplate(l.Model, "")
	if template == "" {
		template = "chatml"
	}
	buf := make([]byte, 4096)
	len := llama.ChatApplyTemplate(template, messages, true, buf)
	result := string(buf[:len])

	// Generate translation using the Llama model with streaming
	response, err := l.generateTextStream(ctx, result, callback)
	if err != nil {
		return model.TranslationResponse{
			Key:            req.Key,
			TargetLanguage: req.TargetLanguage,
			Error:          fmt.Errorf("llama stream generation failed: %v", err),
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

		batch = llama.BatchGetOne([]llama.Token{nextToken})

		// Add to result
		result.WriteString(nextPiece)
	}

	return result.String(), nil
}

// generateTextStream generates text using the Llama model with streaming support
func (l *LlamaTranslator) generateTextStream(ctx context.Context, prompt string, callback func(chunk string) error) (string, error) {
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
	}

	// Generate tokens with streaming
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

		batch = llama.BatchGetOne([]llama.Token{nextToken})

		// Add to result
		result.WriteString(nextPiece)

		// Call callback with the new chunk
		if callback != nil {
			if err := callback(nextPiece); err != nil {
				return result.String(), err
			}
		}
	}

	return result.String(), nil
}

// createTranslationPrompt creates a prompt for translating text
func createTranslationPrompt(req model.TranslationRequest) string {
	// 使用 hunyuan 标识符（Identifier）而不是显示名称，这样模型更容易理解
	targetLanguage := hunyuanLang.MapAppStoreLocaleToHunyuan(req.TargetLanguage)

	// 如果映射失败（返回原值），尝试使用显示名称作为降级
	if targetLanguage == req.TargetLanguage {
		displayName := hunyuanLang.MapAppStoreLocaleToDisplayName(req.TargetLanguage)
		if displayName != req.TargetLanguage {
			targetLanguage = displayName
		} else {
			// 如果都不支持，使用原始代码并发出警告
			log.Warnf("Language '%s' is not supported by Hunyuan model, using raw code", req.TargetLanguage)
			targetLanguage = req.TargetLanguage
		}
	}

	// 创建翻译提示
	prompt := fmt.Sprintf(
		"Translate the following segment into %s, without additional explanation.\n\n%s",
		targetLanguage,
		req.Text,
	)
	log.Infof("prompt: %s", prompt)

	return prompt
}

func cleanLlamaResponse(response string) string {
	// Normalize line endings first - convert \r\n to \n
	response = strings.ReplaceAll(response, "\r\n", "\n")
	
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

	// Remove markdown formatting that the model might have added
	// Remove bold markers **
	trimmed = strings.ReplaceAll(trimmed, "**", "")
	
	// Remove list markers at the start of lines (-, *, bullet)
	paragraphs := strings.Split(trimmed, "\n\n")
	var cleanedParagraphs []string
	
	for _, paragraph := range paragraphs {
		// Process each paragraph - remove leading/trailing whitespace from each line
		lines := strings.Split(paragraph, "\n")
		cleanedLines := make([]string, 0)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Remove markdown list markers at the start
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")
			line = strings.TrimPrefix(line, "• ")
			if line != "" {
				cleanedLines = append(cleanedLines, line)
			}
		}
		// Join lines within paragraph with single newline
		if len(cleanedLines) > 0 {
			cleanedParagraphs = append(cleanedParagraphs, strings.Join(cleanedLines, "\n"))
		}
	}
	
	// Join paragraphs with double newlines to preserve paragraph structure
	trimmed = strings.Join(cleanedParagraphs, "\n\n")

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

// InitLlamaLibrary initializes the llama library with the specified lib path
func InitLlamaLibrary(libPath string) error {
	// Load the llama library
	if err := llama.Load(libPath); err != nil {
		return fmt.Errorf("failed to load llama library from %s: %v", libPath, err)
	}

	// Initialize the llama backend
	llama.Init()

	return nil
}
