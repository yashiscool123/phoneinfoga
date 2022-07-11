package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/output"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"os"
)

var inputNumber string
var pluginPaths []string

func init() {
	// Register command
	rootCmd.AddCommand(scanCmd)

	// Register flags
	scanCmd.PersistentFlags().StringVarP(&inputNumber, "number", "n", "", "The phone number to scan (E164 or international format)")
	// scanCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Text file containing a list of phone numbers to scan (one per line)")
	// scanCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output to save scan results")
	scanCmd.PersistentFlags().StringSliceVar(&pluginPaths, "plugin", []string{}, "Extra scanner plugin to use for the scan")
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a phone number",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runScan()
	},
}

func runScan() {
	fmt.Printf(color.WhiteString("Running scan for phone number %s...\n\n"), inputNumber)

	if valid := number.IsValid(inputNumber); !valid {
		logrus.WithFields(map[string]interface{}{
			"input": inputNumber,
			"valid": valid,
		}).Debug("Input phone number is invalid")
		exitWithError(errors.New("given phone number is not valid"))
	}

	num, err := number.NewNumber(inputNumber)
	if err != nil {
		exitWithError(err)
	}

	remoteLibrary := remote.NewLibrary()
	remote.InitScanners(remoteLibrary)

	for _, pluginPath := range pluginPaths {
		s, err := remote.OpenPlugin(pluginPath)
		if err != nil {
			exitWithError(err)
		}
		remoteLibrary.AddScanner(s)
	}

	result, errs := remoteLibrary.Scan(num)

	err = output.GetOutput(output.Console, os.Stdout).Write(result, errs)
	if err != nil {
		exitWithError(err)
	}
}
