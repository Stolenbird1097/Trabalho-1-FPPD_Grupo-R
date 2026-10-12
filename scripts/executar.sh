#!/usr/bin/env bash
# Executa cada versão N vezes (padrão 5, mínimo exigido) e grava as métricas em CSV.
# Uso:  ./scripts/executar.sh            -> 5 execuções
#       N=10 ./scripts/executar.sh       -> 10 execuções
# Rode SEMPRE na mesma máquina usada no relatório, com o notebook na tomada
# e o mínimo de programas abertos.

set -euo pipefail
cd "$(dirname "$0")/.."

N="${N:-5}"
SAIDA="resultados/resultados.csv"
mkdir -p bin resultados
EXT="$(go env GOEXE)"   # ".exe" no Windows, vazio no Linux/Mac

echo "Compilando..."
for v in v0-sequencial v1-deadlock v2-sincronizacao-excessiva v3-concorrente; do
  go build -o "bin/$v$EXT" "./cmd/$v"
done

echo "versao,total_ms,espera_media_ms,espera_max_ms,max_simultaneos,ciclos_por_filosofo" > "$SAIDA"

for v in v0-sequencial v2-sincronizacao-excessiva v3-concorrente; do
  for i in $(seq 1 "$N"); do
    echo "  $v  execução $i/$N"
    "./bin/$v$EXT" | tail -n 1 >> "$SAIDA"
  done
done

# V1 pode travar: roda com limite de tempo e registra quantas vezes travou.
TRAVOU=0
for i in $(seq 1 "$N"); do
  if timeout 30s ./bin/v1-deadlock$EXT > resultados/v1_saida.txt 2>&1; then
    tail -n 1 resultados/v1_saida.txt >> "$SAIDA"
  else
    TRAVOU=$((TRAVOU + 1))
  fi
done
echo "V1: travou em $TRAVOU de $N execuções" | tee resultados/v1_deadlocks.txt

echo "Resultados em $SAIDA"
