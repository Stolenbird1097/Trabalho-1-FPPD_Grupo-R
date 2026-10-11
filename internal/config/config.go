// Package config centraliza os parâmetros do experimento.
//
// O enunciado exige que TODAS as versões usem os mesmos tempos de pensar e
// comer. Por isso nenhuma versão define tempos próprios: todas chamam
// config.Carregar() e usam os valores daqui.
package config

import (
	"flag"
	"time"
)

const (
	NumFilosofos = 5  // fixo pelo enunciado
	NumCiclos    = 20 // fixo pelo enunciado: 20 ciclos por filósofo
)

type Config struct {
	TempoPensar time.Duration
	TempoComer  time.Duration
	CSV         bool // imprime uma linha CSV no final (usado por scripts/executar.sh)
	Verbose     bool // imprime cada evento (útil para prints e para a evidência do deadlock)
}

// Carregar lê as flags de linha de comando.
// Os valores padrão são os oficiais do grupo: não alterar sem combinar com todos.
func Carregar() Config {
	var c Config
	pensar := flag.Int("pensar", 10, "tempo de pensar em ms")
	comer := flag.Int("comer", 10, "tempo de comer em ms")
	flag.BoolVar(&c.CSV, "csv", false, "imprime o resultado em uma linha CSV")
	flag.BoolVar(&c.Verbose, "v", false, "modo verboso (log de cada evento)")
	flag.Parse()
	c.TempoPensar = time.Duration(*pensar) * time.Millisecond
	c.TempoComer = time.Duration(*comer) * time.Millisecond
	return c
}
