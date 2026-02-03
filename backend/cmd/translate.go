package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tm "github.com/fdddf/opentrans/internal/model"
	"github.com/fdddf/opentrans/internal/translator"
	"github.com/hybridgroup/yzma/pkg/llama"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate xcstrings or plain text files using local Llama model",
	Long: `Translate Localizable.xcstrings or plain text files using local Llama models through the yzma library.
	
This command uses local LLM models for translation, providing privacy-preserving translation
without sending data to external APIs. Requires a local Llama model file.`,
	Run: runTranslate,
}

func init() {
	rootCmd.AddCommand(translateCmd)

	// Required flags
	translateCmd.Flags().String("input-file", "", "Input file path (supports both .xcstrings and plain text files) (required)")
	translateCmd.Flags().String("source-language", "", "Source language code (required)")
	translateCmd.Flags().StringSlice("target-languages", []string{}, "Target language codes (required)")

	// Optional flags
	translateCmd.Flags().Bool("output-with-target-language", false, "Include target language in output filename (e.g., test_chinese.md)")

	// Llama specific flags
	translateCmd.Flags().String("lib-path", "", "Path to the Llama library (required)")
	translateCmd.Flags().String("model-path", "", "Path to the Llama model file (required)")
	translateCmd.Flags().String("grammar-path", "", "Path to the grammar file")
	translateCmd.Flags().Int("threads", 0, "Number of threads to use")
	translateCmd.Flags().Int("seed", -1, "Seed for random number generation (-1 for random)")
	translateCmd.Flags().Int("tokens", 4096, "Maximum number of tokens to generate")
	translateCmd.Flags().Int("top-k", 20, "Top-K sampling parameter")
	translateCmd.Flags().Float64("tfs", 1.0, "Tail free sampling parameter")
	translateCmd.Flags().Float64("top-p", 0.6, "Top-P sampling parameter")
	translateCmd.Flags().Float64("min-p", 0.1, "Minimum P sampling parameter")
	translateCmd.Flags().Float64("typical-p", 1.0, "Typical P sampling parameter")
	translateCmd.Flags().Int("repeat-last-n", 64, "Last n tokens to consider for repetition penalty")
	translateCmd.Flags().Float64("repeat-penalty", 1.05, "Repetition penalty")
	translateCmd.Flags().Float64("frequency-penalty", 0.0, "Frequency penalty")
	translateCmd.Flags().Float64("presence-penalty", 0.0, "Presence penalty")
	translateCmd.Flags().Float64("temperature", 0.7, "Temperature for generation")
	translateCmd.Flags().Bool("verbose-llama", false, "Enable verbose output for Llama operations")

	// Bind flags to Viper
	viper.BindPFlag("global.output_with_target_language", translateCmd.Flags().Lookup("output-with-target-language"))

	// Mark required flags
	translateCmd.MarkFlagRequired("input-file")
	translateCmd.MarkFlagRequired("source-language")
	translateCmd.MarkFlagRequired("target-languages")

	// Bind flags to Viper
	viper.BindPFlag("global.input_file", translateCmd.Flags().Lookup("input-file"))
	viper.BindPFlag("global.source_language", translateCmd.Flags().Lookup("source-language"))
	viper.BindPFlag("global.target_languages", translateCmd.Flags().Lookup("target-languages"))
	viper.BindPFlag("llama.lib_path", translateCmd.Flags().Lookup("lib-path"))
	viper.BindPFlag("llama.model_path", translateCmd.Flags().Lookup("model-path"))
	viper.BindPFlag("llama.grammar_path", translateCmd.Flags().Lookup("grammar-path"))
	viper.BindPFlag("llama.threads", translateCmd.Flags().Lookup("threads"))
	viper.BindPFlag("llama.seed", translateCmd.Flags().Lookup("seed"))
	viper.BindPFlag("llama.tokens", translateCmd.Flags().Lookup("tokens"))
	viper.BindPFlag("llama.top_k", translateCmd.Flags().Lookup("top-k"))
	viper.BindPFlag("llama.tfs", translateCmd.Flags().Lookup("tfs"))
	viper.BindPFlag("llama.top_p", translateCmd.Flags().Lookup("top-p"))
	viper.BindPFlag("llama.minp", translateCmd.Flags().Lookup("min-p"))
	viper.BindPFlag("llama.typical_p", translateCmd.Flags().Lookup("typical-p"))
	viper.BindPFlag("llama.repeat_last_n", translateCmd.Flags().Lookup("repeat-last-n"))
	viper.BindPFlag("llama.repeat_penalty", translateCmd.Flags().Lookup("repeat-penalty"))
	viper.BindPFlag("llama.frequency_penalty", translateCmd.Flags().Lookup("frequency-penalty"))
	viper.BindPFlag("llama.presence_penalty", translateCmd.Flags().Lookup("presence-penalty"))
	viper.BindPFlag("llama.temperature", translateCmd.Flags().Lookup("temperature"))
	viper.BindPFlag("llama.verbose", translateCmd.Flags().Lookup("verbose-llama"))
}

func runTranslate(cmd *cobra.Command, args []string) {
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

	// Get configuration values
	inputFile := viper.GetString("global.input_file")
	if cmd.Flags().Changed("input-file") {
		if val, err := cmd.Flags().GetString("input-file"); err == nil {
			inputFile = val
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

	minP := viper.GetFloat64("llama.min_p")
	if cmd.Flags().Changed("min-p") {
		if val, err := cmd.Flags().GetFloat64("min-p"); err == nil {
			minP = val
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
		fmt.Printf("  Source language: %s\n", sourceLang)
		fmt.Printf("  Target languages: %v\n", targetLangs)
		fmt.Printf("  Model path: %s\n", modelPath)
		fmt.Printf("  Grammar path: %s\n", grammarPath)
		fmt.Printf("  Threads: %d\n", threads)
		fmt.Printf("  Temperature: %.2f\n", temperature)
		fmt.Printf("  Max tokens: %d\n", tokens)
	}

	// Determine file type and handle accordingly
	ext := strings.ToLower(filepath.Ext(inputFile))
	isXCStrings := (ext == ".xcstrings")

	outputWithTargetLanguage := viper.GetBool("global.output_with_target_language")
	if cmd.Flags().Changed("output-with-target-language") {
		if val, err := cmd.Flags().GetBool("output-with-target-language"); err == nil {
			outputWithTargetLanguage = val
		}
	}

	if isXCStrings {
		// Handle xcstrings file
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
			MinP:             minP,
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
		concurrency := 1

		// Create translation service
		service := translator.NewTranslationService(
			provider,
			concurrency,
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

		// Save output file(s)
		if outputWithTargetLanguage {
			// Create a separate output file for each target language
			for _, target := range targetLangs {
				outputFile := strings.TrimSuffix(inputFile, ".xcstrings") + "_" + strings.ReplaceAll(target, "-", "_") + ".xcstrings"

				// Create updated xcstrings object that only contains translations for the current language
				updatedXCStrings := *xcstrings // Create a copy
				// Apply only the translations for the current target language
				var targetResponses []tm.TranslationResponse
				for _, resp := range responses {
					if resp.TargetLanguage == target {
						targetResponses = append(targetResponses, resp)
					}
				}
				translator.ApplyTranslations(&updatedXCStrings, targetResponses)

				if verbose {
					fmt.Printf("Saving output to %s...\n", outputFile)
				}
				err = tm.SaveXCStrings(outputFile, &updatedXCStrings)
				if err != nil {
					fmt.Printf("Error saving output file for %s: %v\n", target, err)
					continue
				}

				fmt.Printf("XCStrings translation to %s completed successfully!\n", target)
				fmt.Printf("Results saved to: %s\n", outputFile)
			}
		} else {
			// Use the original approach with a single output file
			outputFile := strings.TrimSuffix(inputFile, ".xcstrings") + "_translated.xcstrings"

			if verbose {
				fmt.Printf("Saving output to %s...\n", outputFile)
			}
			err = tm.SaveXCStrings(outputFile, xcstrings)
			if err != nil {
				fmt.Printf("Error saving output file: %v\n", err)
				return
			}

			fmt.Printf("XCStrings translation completed successfully!\n")
			fmt.Printf("Results saved to: %s\n", outputFile)
		}
	} else {
		// Handle plain text file
		if verbose {
			fmt.Println("Loading plain text file...")
		}
		content, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading plain text file: %v\n", err)
			return
		}

		textContent := string(content)

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
			MinP:             minP,
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
		concurrency := 1

		// Create translation service
		service := translator.NewTranslationService(
			provider,
			concurrency,
			300*time.Second, // 5 minute timeout (local processing but may take time),
		)

		if verbose {
			fmt.Printf("Translating plain text to %d languages...\n", len(targetLangs))
		}

		ctx := context.Background()

		if outputWithTargetLanguage {
			// Create a separate output file for each target language
			for _, target := range targetLangs {
				// Create a translation request for the plain text content
				req := tm.TranslationRequest{
					Key:            fmt.Sprintf("text_%d", 0),
					Text:           textContent,
					SourceLanguage: sourceLang,
					TargetLanguage: target,
				}

				if verbose {
					fmt.Printf("Translating to %s...\n", target)
				}

				progress := translator.NewVerboseProgressReporter(target, 1, verbose)
				resp, err := service.TranslateBatch(ctx, []tm.TranslationRequest{req}, progress)
				if err != nil {
					fmt.Printf("Translation failed for %s: %v\n", target, err)
					continue
				}

				if len(resp) > 0 && resp[0].Error != nil {
					fmt.Printf("Error translating to %s: %v\n", target, resp[0].Error)
					continue
				}

				if len(resp) > 0 {
					// Set output filename with target language
					outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_" + strings.ReplaceAll(target, "-", "_") + filepath.Ext(inputFile)

					if verbose {
						fmt.Printf("Saving translated text to %s...\n", outputFile)
					}
					err = os.WriteFile(outputFile, []byte(resp[0].TranslatedText), 0644)
					if err != nil {
						fmt.Printf("Error writing output file: %v\n", err)
						continue
					}

					fmt.Printf("Plain text translation to %s completed successfully!\n", target)
					fmt.Printf("Results saved to: %s\n", outputFile)
				}
			}
		} else {
			// Original behavior: translate all languages to a single file
			var outputContent string

			for i, target := range targetLangs {
				// Create a translation request for the plain text content
				req := tm.TranslationRequest{
					Key:            fmt.Sprintf("text_%d", i),
					Text:           textContent,
					SourceLanguage: sourceLang,
					TargetLanguage: target,
				}

				if verbose {
					fmt.Printf("Translating to %s...\n", target)
				}

				progress := translator.NewVerboseProgressReporter(target, 1, verbose)
				resp, err := service.TranslateBatch(ctx, []tm.TranslationRequest{req}, progress)
				if err != nil {
					fmt.Printf("Translation failed for %s: %v\n", target, err)
					continue
				}

				if len(resp) > 0 && resp[0].Error != nil {
					fmt.Printf("Error translating to %s: %v\n", target, resp[0].Error)
					continue
				}

				if len(resp) > 0 {
					if i == 0 {
						outputContent = resp[0].TranslatedText
					} else {
						outputContent += fmt.Sprintf("\n\n--- Translation to %s ---\n%s", target, resp[0].TranslatedText)
					}
				}
			}

			// Set a default output filename for text files
			outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_translated" + filepath.Ext(inputFile)

			if verbose {
				fmt.Printf("Saving translated text to %s...\n", outputFile)
			}
			err = os.WriteFile(outputFile, []byte(outputContent), 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				return
			}

			fmt.Printf("Plain text translation completed successfully!\n")
			fmt.Printf("Results saved to: %s\n", outputFile)
		}
	}
}
