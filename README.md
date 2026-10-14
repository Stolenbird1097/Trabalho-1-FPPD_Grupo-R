# Trabalho 1 – Filósofos Jantando (Concorrência e Sincronização em Go)

Integrantes: Arthur Pimentel · Raffaella Aranha · Tiago Prudencio

Quatro versões do problema dos filósofos jantando, medidas com as mesmas métricas e os mesmos tempos.

| Versão | Pasta | Ideia |
|---|---|---|
| V0 | `codigos/v0-sequencial` | Baseline sem concorrência (base: versão do Moodle) |
| V1 | `codigos/v1-deadlock` | Concorrente, sujeita a deadlock |
| V2 | `codigos/v2-sincronizacao-excessiva` | Sem deadlock, mas limita demais a concorrência |
| V3 | `codigos/v3-concorrente` | Sem deadlock e com mais concorrência |

## Estrutura

```
codigos/             uma pasta por versão (cada uma tem seu main.go)
internal/config      5 filósofos, 20 ciclos, tempos de pensar/comer (iguais para todas)
internal/metrics     as 5 métricas do enunciado + saída em CSV
scripts/             criar_issues.sh (organiza o GitHub) e executar.sh (5 execuções → CSV)
resultados/          CSV, tabela, gráficos e prints
apresentacao/        banco de perguntas e roteiro de execução ao vivo
CHECKLIST.md         tudo que o enunciado cobra, item por item
```

## Como executar

```bash
go version                          # registrar no relatório
go run ./codigos/v3-concorrente         # roda uma versão
go run ./codigos/v1-deadlock -v         # modo verboso (para a evidência do deadlock)
go run -race ./codigos/v2-sincronizacao-excessiva   # checa condições de corrida
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
