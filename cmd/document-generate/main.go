package main

import (
	document "github.com/KinNeko-De/restaurant-document-generate-function/internal/app/document-generate/document"
	operation "github.com/KinNeko-De/restaurant-document-generate-function/internal/app/document-generate/operation"
)

func main() {
	logfile := operation.UseLogFileInGenerated()
	defer operation.CloseLogFile(logfile)

	document.DocumentGenerator{}.GenerateDocument()
}
