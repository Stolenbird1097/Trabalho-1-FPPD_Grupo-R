# Trabalho 1 – Filósofos Jantando (Concorrência e Sincronização em Go)

Integrantes: _Nome A_ · _Nome B_ · _Nome C_

Quatro versões do problema dos filósofos jantando, medidas com as mesmas métricas e os mesmos tempos.

| Versão | Pasta | Ideia |
|---|---|---|
| V0 | `cmd/v0-sequencial` | Baseline sem concorrência (base: versão do Moodle) |
| V1 | `cmd/v1-deadlock` | Concorrente, sujeita a deadlock |
| V2 | `cmd/v2-sincronizacao-excessiva` | Sem deadlock, mas limita demais a concorrência |
| V3 | `cmd/v3-concorrente` | Sem deadlock e com mais concorrência |

## Estrutura

```
cmd/                 uma pasta por versão (cada uma tem seu main.go)
internal/config      5 filósofos, 20 ciclos, tempos de pensar/comer (iguais para todas)
internal/metrics     as 5 métricas do enunciado + saída em CSV
scripts/             criar_issues.sh (organiza o GitHub) e executar.sh (5 execuções → CSV)
resultados/          CSV, tabela, gráficos e prints
relatorio/           ambiente.md, rascunhos e o PDF final
apresentacao/        slides
CHECKLIST.md         tudo que o enunciado cobra, item por item
```

## Como executar

```bash
go version                          # registrar no relatório
go run ./cmd/v3-concorrente         # roda uma versão
go run ./cmd/v1-deadlock -v         # modo verboso (para a evidência do deadlock)
go run -race ./cmd/v2-sincronizacao-excessiva   # checa condições de corrida
./scripts/executar.sh               # 5 execuções de cada versão → resultados/resultados.csv
```

Flags comuns: `-pensar 10 -comer 10` (ms), `-csv`, `-v`. **Os valores padrão são os oficiais; não mudar nas medições finais.**

## Como cada versão usa as métricas

```go
cfg := config.Carregar()
m := metrics.Nova("V3", config.NumFilosofos)
m.Iniciar()
// cada filósofo, em cada um dos config.NumCiclos ciclos:
//   time.Sleep(cfg.TempoPensar)
//   inicio := time.Now(); <pega os dois garfos>; m.RegistrarEspera(id, time.Since(inicio))
//   m.ComecouAComer(); time.Sleep(cfg.TempoComer); m.TerminouDeComer()
//   <solta os garfos>; m.CicloConcluido(id)
m.Finalizar()
if cfg.CSV { fmt.Println(m.LinhaCSV()) } else { m.Imprimir() }
```

## Divisão do trabalho

| | Pessoa A | Pessoa B | Pessoa C |
|---|---|---|---|
| Código | Setup, config/métricas, **V0**, **V2** | **V1**, **V3** | Revisão cruzada, CI |
| Experimentos | Investigação da V2 | Evidência do deadlock da V1 | Script, 5 execuções, tabela, gráficos |
| Relatório | Perguntas 1–3 de V0/V2 | Perguntas 1–6 de V1/V3 | Ambiente, págs. 1 e 2, perguntas 7–12 |
| Apresentação | Fala sobre V0/V2 | Demo ao vivo | Slides |

Conclusão, revisão final do PDF e ensaio são de todos. Todos precisam saber explicar todas as versões (cai na Prova 1).

## Fluxo de trabalho

1. Pegue uma issue e crie uma branch: `git switch -c v1-deadlock`
2. Commits pequenos; antes do push: `gofmt -w . && go vet ./...`
3. Abra um Pull Request com `Closes #N`; outra pessoa revisa e aprova
4. Merge na `main` → a issue fecha sozinha e o quadro se atualiza

## Ambiente de execução

Ver `relatorio/ambiente.md`.
