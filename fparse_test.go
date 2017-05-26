package main

import "testing"
import "encoding/json"

func TestFparse(t *testing.T) {
	var parser FParse

	rawConfig := `{
		"TipoLinha": "linha cabecalho",
		""
	}`

	configJSON := json.NewDecoder().Decode()
	parser = &FParse{}

	linhas := parser.Analisar("./teste.txt", &configJSON)

}
