package main

import (
	"encoding/base64"
	"testing"

	"gopkg.in/yaml.v2"

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

func TestDataToBase64(t *testing.T) {
	input := `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: my-namespace
type: Opaque
data:
  foo: bar
  secret: |
    hello: world
    answer: 42`
	output, err := dataToBase64(input)
	assert.NilError(t, err)

	outputYaml := struct {
		ApiVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
		Type string `yaml:"type"`
		Data struct {
			Foo    string `yaml:"foo"`
			Secret string `yaml:"secret"`
		} `yaml:"data"`
	}{}
	err = yaml.Unmarshal([]byte(output), &outputYaml)
	assert.NilError(t, err)

	assert.Equal(t, outputYaml.ApiVersion, "v1")
	assert.Equal(t, outputYaml.Kind, "Secret")
	assert.Equal(t, outputYaml.Metadata.Name, "my-secret")
	assert.Equal(t, outputYaml.Metadata.Namespace, "my-namespace")
	assert.Equal(t, outputYaml.Type, "Opaque")

	fooDecoded, err := base64.StdEncoding.DecodeString(outputYaml.Data.Foo)
	assert.Equal(t, string(fooDecoded), "bar")

	barDecoded, err := base64.StdEncoding.DecodeString(outputYaml.Data.Secret)
	assert.Equal(t, string(barDecoded), "hello: world\nanswer: 42")
}
