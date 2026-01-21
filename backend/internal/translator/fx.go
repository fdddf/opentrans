package translator

import (
	"go.uber.org/fx"
)

// TranslatorModule is the FX module for translator
// Note: The actual TranslationService is defined in translator.go
// This module can be used to provide additional translator-related dependencies if needed
var TranslatorModule = fx.Module("translator")