package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brianstarke/schemactor/internal/consolidator"
	"github.com/brianstarke/schemactor/internal/verifier"
)

// Version is injected at build time via ldflags
var Version = "dev"

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorBold   = "\033[1m"
)

func main() {
	// Default directories
	inputDir := "./sample_migrations"
	outputDir := "./output"
	verify := false

	// Parse command line arguments
	args := []string{}
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-h" || arg == "--help" {
			printUsage()
			os.Exit(0)
		} else if arg == "--version" || arg == "-V" {
			printVersion()
			os.Exit(0)
		} else if arg == "-v" || arg == "--verify" {
			verify = true
		} else {
			args = append(args, arg)
		}
	}

	// Parse positional arguments
	if len(args) > 0 {
		inputDir = args[0]
	}
	if len(args) > 1 {
		outputDir = args[1]
	}

	// Verify input directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		printError(fmt.Sprintf("Input directory does not exist: %s", inputDir))
		os.Exit(1)
	}

	// Count input migrations
	files, err := filepath.Glob(filepath.Join(inputDir, "*.up.sql"))
	if err != nil {
		printError(fmt.Sprintf("Error reading input directory: %v", err))
		os.Exit(1)
	}
	inputCount := len(files)

	if inputCount == 0 {
		printError(fmt.Sprintf("No migration files found in: %s", inputDir))
		os.Exit(1)
	}

	// Print header
	printHeader()
	fmt.Printf("Input:  %s%s%s\n", colorCyan, inputDir, colorReset)
	fmt.Printf("Output: %s%s%s\n", colorCyan, outputDir, colorReset)
	fmt.Println()

	// Run consolidation
	c := consolidator.NewConsolidator(inputDir, outputDir, false)
	if err := c.Consolidate(false); err != nil {
		printError(fmt.Sprintf("Consolidation failed: %v", err))
		os.Exit(1)
	}

	// Count output migrations
	files, err = filepath.Glob(filepath.Join(outputDir, "*.up.sql"))
	if err != nil {
		printError(fmt.Sprintf("Error reading output directory: %v", err))
		os.Exit(1)
	}
	outputCount := len(files)

	// Print success
	printSuccess(inputCount, outputCount)

	// Run verification if requested
	if verify {
		fmt.Println()
		fmt.Printf("%s%sVerifying migrations...%s\n", colorBold, colorYellow, colorReset)

		v := verifier.NewVerifier(outputDir, true)
		ctx := context.Background()

		if err := v.Verify(ctx); err != nil {
			printError(fmt.Sprintf("Verification failed: %v", err))
			os.Exit(1)
		}

		fmt.Println()
		fmt.Printf("%s✓ Verification successful!%s All migrations can be applied and rolled back.\n",
			colorGreen+colorBold, colorReset)
		fmt.Println()
	}
}

func printVersion() {
	fmt.Printf("%s%s%s version %s%s%s\n", colorBold, colorCyan, "schemactor", colorGreen, Version, colorReset)
}

func printUsage() {
	fmt.Printf("\n%s%sSCHEMACTOR%s - SQL Migration Consolidator\n", colorBold, colorCyan, colorReset)
	fmt.Printf("\n%sUsage:%s\n", colorBold, colorReset)
	fmt.Printf("  schemactor [options] [input_dir] [output_dir]\n")
	fmt.Printf("\n%sOptions:%s\n", colorBold, colorReset)
	fmt.Printf("  %s-V, --version%s  Show version information\n", colorYellow, colorReset)
	fmt.Printf("  %s-v, --verify%s   Verify consolidated migrations with PostgreSQL (requires Docker)\n", colorYellow, colorReset)
	fmt.Printf("  %s-h, --help%s     Show this help message\n", colorYellow, colorReset)
	fmt.Printf("\n%sArguments:%s\n", colorBold, colorReset)
	fmt.Printf("  %sinput_dir%s   Directory containing migration files (default: ./sample_migrations)\n", colorYellow, colorReset)
	fmt.Printf("  %soutput_dir%s  Directory for consolidated migrations (default: ./output)\n", colorYellow, colorReset)
	fmt.Printf("\n%sExamples:%s\n", colorBold, colorReset)
	fmt.Printf("  %sschemactor%s\n", colorGray, colorReset)
	fmt.Printf("  %sschemactor --version%s\n", colorGray, colorReset)
	fmt.Printf("  %sschemactor --verify%s\n", colorGray, colorReset)
	fmt.Printf("  %sschemactor ./migrations ./consolidated%s\n", colorGray, colorReset)
	fmt.Printf("  %sschemactor --verify ./migrations ./consolidated%s\n", colorGray, colorReset)
	fmt.Println()
}

func printHeader() {
	fmt.Println()
	fmt.Printf("%s%sSCHEMACTOR%s Consolidator\n", colorBold, colorPurple, colorReset)
	fmt.Println()
}

func printSuccess(inputCount, outputCount int) {
	reduction := 0
	if inputCount > 0 {
		reduction = (inputCount - outputCount) * 100 / inputCount
	}

	fmt.Println()
	fmt.Printf("%s✓ Complete!%s\n", colorGreen+colorBold, colorReset)
	fmt.Println()
	fmt.Printf("Input migrations:  %s%d%s\n", colorCyan, inputCount, colorReset)
	fmt.Printf("Output migrations: %s%d%s\n", colorCyan, outputCount, colorReset)
	fmt.Printf("Reduction:         %s%d%%%s (%s%d%s eliminated)\n",
		colorGreen+colorBold, reduction, colorReset,
		colorYellow, inputCount-outputCount, colorReset)
	fmt.Println()
}

func printError(msg string) {
	fmt.Println()
	fmt.Printf("%s✗ Error:%s %s\n", colorRed+colorBold, colorReset, msg)
	fmt.Println()
}
