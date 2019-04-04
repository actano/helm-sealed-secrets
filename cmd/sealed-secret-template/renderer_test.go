package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetInputOutputPaths(t *testing.T) {
	matches := []string{
		"./secret/subpath1/secret1.template.yaml",
		"./secret/subpath2/secret2.template.yaml",
	}
	inputDir := "./secret"
	outputDir := "./sealed"
	inputOutputPaths, err := GetInputOutputPaths(matches, inputDir, outputDir)

	assert.NilError(t, err)
	assert.DeepEqual(t, inputOutputPaths, []InputOutputPaths{
		{
			InputPath:  "./secret/subpath1/secret1.template.yaml",
			OutputPath: "sealed/subpath1/secret1.sealed.yaml",
		},
		{
			InputPath:  "./secret/subpath2/secret2.template.yaml",
			OutputPath: "sealed/subpath2/secret2.sealed.yaml",
		},
	})
}
