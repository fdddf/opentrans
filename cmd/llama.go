package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	tm "github.com/fdddf/xcstrings-translator/internal/model"
	"github.com/fdddf/xcstrings-translator/internal/translator"
	"github.com/hybridgroup/yzma/pkg/llama"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var llamaCmd = &cobra.Command{
	Use:   "llama",
	Short: "Translate xcstrings using local Llama model",
	Long: `Translate Localizable.xcstrings file using local Llama models through the yzma library.
	
This translator uses local LLM models for translation, providing privacy-preserving translation
without sending data to external APIs. Requires a local Llama model file.`,
	Run: runLlamaTranslate,
}

func init() {
	rootCmd.AddCommand(llamaCmd)

	// Llama specific flags
	llamaCmd.Flags().String("lib-path", "", "Path to the Llama library (required)")
	llamaCmd.Flags().String("model-path", "", "Path to the Llama model file (required)")
	llamaCmd.Flags().String("grammar-path", "", "Path to the grammar file")
	llamaCmd.Flags().Int("threads", 0, "Number of threads to use")
	llamaCmd.Flags().Int("seed", -1, "Seed for random number generation (-1 for random)")
	llamaCmd.Flags().Int("tokens", 4096, "Maximum number of tokens to generate")
	llamaCmd.Flags().Int("top-k", 20, "Top-K sampling parameter")
	llamaCmd.Flags().Float64("tfs", 1.0, "Tail free sampling parameter")
	llamaCmd.Flags().Float64("top-p", 0.6, "Top-P sampling parameter")
	llamaCmd.Flags().Float64("typical-p", 1.0, "Typical P sampling parameter")
	llamaCmd.Flags().Int("repeat-last-n", 64, "Last n tokens to consider for repetition penalty")
	llamaCmd.Flags().Float64("repeat-penalty", 1.05, "Repetition penalty")
	llamaCmd.Flags().Float64("frequency-penalty", 0.0, "Frequency penalty")
	llamaCmd.Flags().Float64("presence-penalty", 0.0, "Presence penalty")
	llamaCmd.Flags().Float64("temperature", 0.7, "Temperature for generation")
	llamaCmd.Flags().Bool("verbose-llama", false, "Enable verbose output for Llama operations")

	// Bind flags to Viper
	viper.BindPFlag("llama.lib_path", llamaCmd.Flags().Lookup("lib-path"))
	viper.BindPFlag("llama.model_path", llamaCmd.Flags().Lookup("model-path"))
	viper.BindPFlag("llama.grammar_path", llamaCmd.Flags().Lookup("grammar-path"))
	viper.BindPFlag("llama.threads", llamaCmd.Flags().Lookup("threads"))
	viper.BindPFlag("llama.seed", llamaCmd.Flags().Lookup("seed"))
	viper.BindPFlag("llama.tokens", llamaCmd.Flags().Lookup("tokens"))
	viper.BindPFlag("llama.top_k", llamaCmd.Flags().Lookup("top-k"))
	viper.BindPFlag("llama.tfs", llamaCmd.Flags().Lookup("tfs"))
	viper.BindPFlag("llama.top_p", llamaCmd.Flags().Lookup("top-p"))
	viper.BindPFlag("llama.typical_p", llamaCmd.Flags().Lookup("typical-p"))
	viper.BindPFlag("llama.repeat_last_n", llamaCmd.Flags().Lookup("repeat-last-n"))
	viper.BindPFlag("llama.repeat_penalty", llamaCmd.Flags().Lookup("repeat-penalty"))
	viper.BindPFlag("llama.frequency_penalty", llamaCmd.Flags().Lookup("frequency-penalty"))
	viper.BindPFlag("llama.presence_penalty", llamaCmd.Flags().Lookup("presence-penalty"))
	viper.BindPFlag("llama.temperature", llamaCmd.Flags().Lookup("temperature"))
	viper.BindPFlag("llama.verbose", llamaCmd.Flags().Lookup("verbose-llama"))
}

func runLlamaTranslate(cmd *cobra.Command, args []string) {
	// Get Llama specific config first to check if model path is provided
	modelPath := viper.GetString("llama.model_path")
	if cmd.Flags().Changed("model-path") {
		if val, err := cmd.Flags().GetString("model-path"); err == nil {
			modelPath = val
		}
	}

	// Validate required parameters
	if modelPath == "" {
		fmt.Println("Error: model-path is required for Llama translator")
		return
	}

	libPath := viper.GetString("llama.lib_path")
	if libPath == "" {
		fmt.Println("Error: lib-path is required for Llama translator")
		return
	}

	// Initialize the llama library first
	// Try to load without a specific path first (will use default path or environment variable)
	if err := llama.Load(libPath); err != nil {
		// If that fails, the user may need to set the YZMA_LIB environment variable
		// or provide the path to the llama library in some other way
		fmt.Println("Warning: unable to load llama library automatically:", err.Error())
		fmt.Println("Make sure the llama library is installed and accessible.")
		fmt.Println("You may need to set the YZMA_LIB environment variable to point to the library path.")
		return
	}

	// Initialize the llama backend
	llama.Init()
	defer llama.Close()

	if !viper.GetBool("llama.verbose") {
		llama.LogSet(llama.LogSilent())
	}

	// Get configuration values with fallbacks
	inputFile := viper.GetString("global.input_file")
	if cmd.Flags().Changed("input") {
		if val, err := cmd.Flags().GetString("input"); err == nil {
			inputFile = val
		}
	}

	outputFile := viper.GetString("global.output_file")
	if cmd.Flags().Changed("output") {
		if val, err := cmd.Flags().GetString("output"); err == nil {
			outputFile = val
		}
	}

	sourceLang := viper.GetString("global.source_language")
	if cmd.Flags().Changed("source-language") {
		if val, err := cmd.Flags().GetString("source-language"); err == nil {
			sourceLang = val
		}
	}
	sourceLang = normalizeLanguageCode(sourceLang)

	targetLangs := viper.GetStringSlice("global.target_languages")
	if cmd.Flags().Changed("target-languages") {
		if val, err := cmd.Flags().GetStringSlice("target-languages"); err == nil {
			targetLangs = val
		}
	}
	for i, lang := range targetLangs {
		targetLangs[i] = normalizeLanguageCode(lang)
	}

	concurrency := viper.GetInt("global.concurrency")
	if cmd.Flags().Changed("concurrency") {
		if val, err := cmd.Flags().GetInt("concurrency"); err == nil {
			concurrency = val
		}
	}

	verbose := viper.GetBool("global.verbose")
	if cmd.Flags().Changed("verbose") {
		if val, err := cmd.Flags().GetBool("verbose"); err == nil {
			verbose = val
		}
	}

	// Get Llama specific config
	// modelPath is already obtained above
	grammarPath := viper.GetString("llama.grammar_path")
	if cmd.Flags().Changed("grammar-path") {
		if val, err := cmd.Flags().GetString("grammar-path"); err == nil {
			grammarPath = val
		}
	}

	threads := viper.GetInt("llama.threads")
	if cmd.Flags().Changed("threads") {
		if val, err := cmd.Flags().GetInt("threads"); err == nil {
			threads = val
		}
	}

	seed := viper.GetInt("llama.seed")
	if cmd.Flags().Changed("seed") {
		if val, err := cmd.Flags().GetInt("seed"); err == nil {
			seed = val
		}
	}

	tokens := viper.GetInt("llama.tokens")
	if cmd.Flags().Changed("tokens") {
		if val, err := cmd.Flags().GetInt("tokens"); err == nil {
			tokens = val
		}
	}

	topK := viper.GetInt("llama.top_k")
	if cmd.Flags().Changed("top-k") {
		if val, err := cmd.Flags().GetInt("top-k"); err == nil {
			topK = val
		}
	}

	tfs := viper.GetFloat64("llama.tfs")
	if cmd.Flags().Changed("tfs") {
		if val, err := cmd.Flags().GetFloat64("tfs"); err == nil {
			tfs = val
		}
	}

	topP := viper.GetFloat64("llama.top_p")
	if cmd.Flags().Changed("top-p") {
		if val, err := cmd.Flags().GetFloat64("top-p"); err == nil {
			topP = val
		}
	}

	typicalP := viper.GetFloat64("llama.typical_p")
	if cmd.Flags().Changed("typical-p") {
		if val, err := cmd.Flags().GetFloat64("typical-p"); err == nil {
			typicalP = val
		}
	}

	repeatLastN := viper.GetInt("llama.repeat_last_n")
	if cmd.Flags().Changed("repeat-last-n") {
		if val, err := cmd.Flags().GetInt("repeat-last-n"); err == nil {
			repeatLastN = val
		}
	}

	repeatPenalty := viper.GetFloat64("llama.repeat_penalty")
	if cmd.Flags().Changed("repeat-penalty") {
		if val, err := cmd.Flags().GetFloat64("repeat-penalty"); err == nil {
			repeatPenalty = val
		}
	}

	frequencyPenalty := viper.GetFloat64("llama.frequency_penalty")
	if cmd.Flags().Changed("frequency-penalty") {
		if val, err := cmd.Flags().GetFloat64("frequency-penalty"); err == nil {
			frequencyPenalty = val
		}
	}

	presencePenalty := viper.GetFloat64("llama.presence_penalty")
	if cmd.Flags().Changed("presence-penalty") {
		if val, err := cmd.Flags().GetFloat64("presence-penalty"); err == nil {
			presencePenalty = val
		}
	}

	temperature := viper.GetFloat64("llama.temperature")
	if cmd.Flags().Changed("temperature") {
		if val, err := cmd.Flags().GetFloat64("temperature"); err == nil {
			temperature = val
		}
	}

	llamaVerbose := viper.GetBool("llama.verbose")
	if cmd.Flags().Changed("verbose-llama") {
		if val, err := cmd.Flags().GetBool("verbose-llama"); err == nil {
			llamaVerbose = val
		}
	}

	if verbose {
		fmt.Printf("Starting Llama Translate with:\n")
		fmt.Printf("  Input file: %s\n", inputFile)
		fmt.Printf("  Output file: %s\n", outputFile)
		fmt.Printf("  Source language: %s\n", sourceLang)
		fmt.Printf("  Target languages: %v\n", targetLangs)
		fmt.Printf("  Concurrency: %d\n", concurrency)
		fmt.Printf("  Model path: %s\n", modelPath)
		fmt.Printf("  Grammar path: %s\n", grammarPath)
		fmt.Printf("  Threads: %d\n", threads)
		fmt.Printf("  Temperature: %.2f\n", temperature)
		fmt.Printf("  Max tokens: %d\n", tokens)
	}

	// Load xcstrings file
	if verbose {
		fmt.Println("Loading xcstrings file...")
	}
	xcstrings, err := tm.LoadXCStrings(inputFile)
	if err != nil {
		fmt.Printf("Error loading xcstrings file: %v\n", err)
		return
	}

	// Override source language if specified
	if sourceLang != "" {
		xcstrings.SourceLanguage = sourceLang
	}

	// Create translator options
	options := translator.LlamaOptions{
		ModelPath:        modelPath,
		GrammarPath:      grammarPath,
		Threads:          threads,
		Seed:             seed,
		Tokens:           tokens,
		TopK:             topK,
		Tfs:              tfs,
		TopP:             topP,
		TypicalP:         typicalP,
		RepeatLastN:      repeatLastN,
		RepeatPenalty:    repeatPenalty,
		FrequencyPenalty: frequencyPenalty,
		PresencePenalty:  presencePenalty,
		Temperature:      temperature,
		Verbose:          llamaVerbose,
	}

	// Create translator
	provider, err := translator.NewLlamaTranslator(options)
	if err != nil {
		fmt.Printf("Error creating Llama translator: %v\n", err)
		return
	}
	defer provider.Close()

	// For llama models, force concurrency to 1 due to thread safety issues with llama library
	effectiveConcurrency := concurrency
	if effectiveConcurrency > 1 {
		fmt.Println("Warning: Llama models require single-threaded access. Concurrency limited to 1.")
		effectiveConcurrency = 1
	}

	// Create translation service
	service := translator.NewTranslationService(
		provider,
		effectiveConcurrency,
		300*time.Second, // 5 minute timeout (local processing but may take time),
	)

	// Run translation
	if verbose {
		fmt.Println("Starting translation...")
	}
	ctx := context.Background()
	var responses []tm.TranslationResponse
	for _, target := range targetLangs {
		reqs := translator.CreateTranslationRequestsForLanguage(xcstrings, target)
		if len(reqs) == 0 {
			continue
		}

		if verbose {
			fmt.Printf("Translating to %s (%d strings)...\n", target, len(reqs))
		}

		progress := translator.NewVerboseProgressReporter(target, len(reqs), verbose)
		batchResponses, err := service.TranslateBatch(ctx, reqs, progress)
		responses = append(responses, batchResponses...)
		if err != nil {
			fmt.Printf("Translation failed for %s: %v\n", target, err)
			return
		}
	}

	if len(responses) == 0 {
		fmt.Println("No strings to translate. Exiting.")
		return
	}

	// Process results
	successCount := 0
	errorCount := 0
	for _, resp := range responses {
		if resp.Error != nil {
			if verbose {
				fmt.Printf("Error translating %s to %s: %v\n", resp.Key, resp.TargetLanguage, resp.Error)
			}
			errorCount++
		} else {
			successCount++
		}
	}

	if verbose {
		fmt.Printf("Translation completed: %d successful, %d failed\n", successCount, errorCount)
	}

	if errorCount > 0 {
		fmt.Println("Errors detected during translation. Stopping without applying translations.")
		return
	}

	// Apply translations
	if verbose {
		fmt.Println("Applying translations...")
	}
	translator.ApplyTranslations(xcstrings, responses)

	// Save output
	if verbose {
		fmt.Printf("Saving output to %s...\n", outputFile)
	}
	err = tm.SaveXCStrings(outputFile, xcstrings)
	if err != nil {
		fmt.Printf("Error saving output file: %v\n", err)
		return
	}

	fmt.Printf("Translation completed successfully!\n")
	fmt.Printf("Results saved to: %s\n", outputFile)
}

func normalizeLanguageCode(language string) string {
	if strings.TrimSpace(language) == "" {
		return language
	}

	normalized := strings.ToLower(strings.TrimSpace(language))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.Join(strings.Fields(normalized), " ")
	normalized = strings.ReplaceAll(normalized, "(", "")
	normalized = strings.ReplaceAll(normalized, ")", "")

	codeByToken := map[string]string{
		"ar":      "ar",
		"cs":      "cs",
		"da":      "da",
		"de":      "de",
		"en":      "en",
		"es":      "es",
		"fi":      "fi",
		"fr":      "fr",
		"he":      "he",
		"hi":      "hi",
		"id":      "id",
		"it":      "it",
		"ja":      "ja",
		"ko":      "ko",
		"nl":      "nl",
		"no":      "no",
		"pl":      "pl",
		"pt":      "pt",
		"ru":      "ru",
		"sv":      "sv",
		"th":      "th",
		"tr":      "tr",
		"vi":      "vi",
		"zh-hans": "zh-Hans",
		"zh-hant": "zh-Hant",
	}
	if code, ok := codeByToken[normalized]; ok {
		return code
	}

	nameByToken := map[string]string{
		"arabic":               "ar",
		"chinese simplified":   "zh-Hans",
		"chinese traditional":  "zh-Hant",
		"czech":                "cs",
		"danish":               "da",
		"dutch":                "nl",
		"english":              "en",
		"finnish":              "fi",
		"french":               "fr",
		"german":               "de",
		"hebrew":               "he",
		"hindi":                "hi",
		"indonesian":           "id",
		"italian":              "it",
		"japanese":             "ja",
		"korean":               "ko",
		"norwegian":            "no",
		"polish":               "pl",
		"portuguese":           "pt",
		"russian":              "ru",
		"simplified chinese":   "zh-Hans",
		"spanish":              "es",
		"swedish":              "sv",
		"thai":                 "th",
		"traditional chinese":  "zh-Hant",
		"turkish":              "tr",
		"vietnamese":           "vi",
	}
	if code, ok := nameByToken[normalized]; ok {
		return code
	}

	return language
}
