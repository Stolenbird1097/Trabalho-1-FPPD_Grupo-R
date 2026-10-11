#!/usr/bin/env bash
# Cria labels, milestones e TODAS as issues do Trabalho 1 no GitHub.
#
# Pré-requisitos:
#   1. GitHub CLI instalado e logado:  gh auth login
#   2. Rodar de dentro do repositório já criado no GitHub (git remote configurado)
#   3. Preencher os usuários do GitHub de cada integrante:
#
#      PESSOA_A=usuarioA PESSOA_B=usuarioB PESSOA_C=usuarioC ./scripts/criar_issues.sh
#
#   (Os outros dois precisam já ter aceitado o convite de colaborador,
#    senão o GitHub não deixa atribuir a issue a eles.)
#
# Rode UMA vez só; rodar de novo duplica as issues.

set -euo pipefail

A="${PESSOA_A:-}"; B="${PESSOA_B:-}"; C="${PESSOA_C:-}"
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
echo "Repositório: $REPO"

# ---------- Labels ----------
label() { gh label create "$1" --color "$2" --description "$3" --force >/dev/null; }
label "v0"              "c5def5" "Versão sequencial (baseline)"
label "v1"              "d93f0b" "Versão com deadlock"
label "v2"              "fbca04" "Sem deadlock, sincronização excessiva"
label "v3"              "0e8a16" "Versão concorrente aprimorada"
label "codigo"          "1d76db" "Implementação em Go"
label "experimento"     "5319e7" "Medições e execuções"
label "relatorio"       "0052cc" "Relatório em PDF"
label "apresentacao"    "b60205" "Apresentação presencial"
label "infra"           "bfdadc" "Repositório, scripts, CI"
label "obrigatorio"     "e11d21" "Item cobrado explicitamente no enunciado"

# ---------- Milestones ----------
milestone() {
  gh api "repos/$REPO/milestones" -f title="$1" -f description="$2" >/dev/null 2>&1 || true
}
M1="1. Setup";            milestone "$M1" "Repositório, código compartilhado e ambiente"
M2="2. Implementações";   milestone "$M2" "V0, V1, V2 e V3 funcionando"
M3="3. Experimentos";     milestone "$M3" "5 execuções por versão, tabela e gráficos"
M4="4. Relatório";        milestone "$M4" "PDF completo respondendo as 12 perguntas"
M5="5. Apresentação";     milestone "$M5" "Slides, demo e ensaio"

# ---------- Issues ----------
issue() { # titulo labels milestone responsavel corpo
  local args=(--title "$1" --label "$2" --milestone "$3" --body "$5")
  [[ -n "$4" ]] && args+=(--assignee "$4")
  gh issue create "${args[@]}" >/dev/null
  echo "  ✔ $1"
}

echo "Criando issues..."

# ===== 1. Setup =====
issue "Configurar repositório e regras de trabalho" "infra" "$M1" "$A" "$(cat <<'EOF'
- [ ] Criar o repositório e convidar os outros dois integrantes como colaboradores
- [ ] Trocar `SEU-USUARIO` no `go.mod` pelo dono real do repositório
- [ ] Proteger a branch `main` (merge só via Pull Request com 1 aprovação)
- [ ] Combinar no README: uma branch por issue (`v1-deadlock`, `relatorio-q7`...), PR fecha a issue com `Closes #N`
- [ ] Rodar `scripts/criar_issues.sh`
EOF
)"

issue "Fixar parâmetros e métricas comuns a todas as versões" "infra,codigo,obrigatorio" "$M1" "$A" "$(cat <<'EOF'
Enunciado: *"todas as versões deverão utilizar os mesmos tempos para as operações simuladas de pensar e comer"*.

- [ ] Revisar `internal/config` (5 filósofos, 20 ciclos, tempos de pensar/comer) e o grupo **aprovar os tempos oficiais**
- [ ] Revisar `internal/metrics` (tempo total, espera média, espera máxima, máx. simultâneos, ciclos por filósofo)
- [ ] Escrever no README como cada versão deve chamar as métricas
- [ ] Ler o link de medição de tempo em Go disponibilizado no Moodle e confirmar que o método bate com o nosso
EOF
)"

issue "Documentar o ambiente de execução" "relatorio,obrigatorio" "$M1" "$C" "$(cat <<'EOF'
O relatório precisa informar (seção "Ambiente de execução"):
- [ ] Linguagem e versão (`go version`)
- [ ] Sistema operacional
- [ ] Processador (modelo)
- [ ] Número de cores disponíveis (`runtime.NumCPU()` e o que o SO informa)
- [ ] Comandos de compilação e execução
- [ ] Garantir que a máquina usada é multicore e que **todas as medições finais rodam na mesma máquina**
Registrar em `relatorio/ambiente.md`.
EOF
)"

# ===== 2. Implementações =====
issue "V0 – Versão sequencial (baseline)" "v0,codigo,obrigatorio" "$M2" "$A" "$(cat <<'EOF'
- [ ] Partir da versão disponibilizada no Moodle (no card do trabalho)
- [ ] Filósofos executam um após o outro, sem concorrência
- [ ] 20 ciclos: pensar → adquirir → comer → liberar
- [ ] Usar `internal/config` e `internal/metrics`
- [ ] Código em `cmd/v0-sequencial/main.go`, formatado com `gofmt`
- [ ] Comentário no topo citando a fonte original e o que foi adaptado
EOF
)"

issue "V1 – Versão concorrente sujeita a deadlock" "v1,codigo,obrigatorio" "$M2" "$B" "$(cat <<'EOF'
- [ ] Buscar uma solução em Go em que cada filósofo pega o garfo da esquerda e depois o da direita (ex.: `sync.Mutex` por garfo)
- [ ] Citar a fonte no comentário do topo
- [ ] Usar `internal/config` e `internal/metrics`
- [ ] Modo verboso (`-v`) que mostra "Fx pegou garfo y / Fx esperando garfo z" para ficar evidente na demonstração
- [ ] Opção para **forçar** o deadlock de forma reproduzível (ex.: pequena pausa entre pegar o 1º e o 2º garfo), para não depender de sorte na apresentação
- [ ] Código em `cmd/v1-deadlock/main.go`
EOF
)"

issue "V1 – Evidência e análise do deadlock" "v1,experimento,relatorio,obrigatorio" "$M3" "$B" "$(cat <<'EOF'
A V1 deve ser analisada separadamente, com evidências do bloqueio.

- [ ] Executar até travar e capturar tela (Go mostra `fatal error: all goroutines are asleep - deadlock!` quando todas param; se não aparecer, usar `kill -QUIT <pid>` / Ctrl+\ para o dump das goroutines)
- [ ] Salvar prints em `resultados/prints/v1-*.png`
- [ ] Tabela: **quais recursos cada filósofo possui** e **por qual recurso cada um espera**
- [ ] Explicar a **espera circular** (F1→F2→…→F5→F1)
- [ ] Mostrar as **4 condições de Coffman** presentes: exclusão mútua, segura-e-espera, ausência de preempção, espera circular
- [ ] Registrar quantas execuções travaram e quantas (se alguma) terminaram
EOF
)"

issue "V2 – Sem deadlock com sincronização excessiva" "v2,codigo,obrigatorio" "$M2" "$A" "$(cat <<'EOF'
- [ ] Buscar solução que evita deadlock limitando demais a concorrência (ex.: um único mutex global para a mesa, "garçom" que deixa só 1 comer por vez)
- [ ] Citar a fonte; se o mecanismo não for nativo em Go, adaptar/refatorar e explicar
- [ ] Acesso aos recursos deve estar **correto** (rodar com `go run -race`)
- [ ] Usar `internal/config` e `internal/metrics`
- [ ] Código em `cmd/v2-sincronizacao-excessiva/main.go`
EOF
)"

issue "V3 – Versão concorrente aprimorada" "v3,codigo,obrigatorio" "$M2" "$B" "$(cat <<'EOF'
- [ ] Buscar solução sem deadlock com mais concorrência: ordenação na aquisição dos garfos (pegar sempre o de menor número primeiro), canais, garçom com semáforo de 4 vagas etc.
- [ ] Justificar a escolha da estratégia
- [ ] Rodar com `go run -race`
- [ ] Usar `internal/config` e `internal/metrics`
- [ ] Código em `cmd/v3-concorrente/main.go`
- [ ] Explicar **por que evita deadlock** (qual condição de Coffman é quebrada)
- [ ] Discutir **possibilidade de starvation**, usando os ciclos por filósofo e a espera máxima como evidência
EOF
)"

issue "Revisão cruzada do código e CI" "infra,codigo" "$M2" "$C" "$(cat <<'EOF'
- [ ] Revisar os PRs de V0–V3 (cada integrante revisa código que não escreveu)
- [ ] Conferir que todas as versões usam os mesmos tempos e as mesmas métricas
- [ ] CI verde (`gofmt`, `go vet`, `go build`)
- [ ] V2 e V3 sem erros no `-race`
EOF
)"

# ===== 3. Experimentos =====
issue "Script de execução e coleta em CSV" "infra,experimento" "$M3" "$C" "$(cat <<'EOF'
- [ ] Revisar `scripts/executar.sh` (5 execuções de V0, V2 e V3 → `resultados/resultados.csv`)
- [ ] V1 roda com timeout para não travar o script
- [ ] Script `scripts/resumo.py` (ou planilha) que calcula média e desvio de cada métrica por versão
EOF
)"

issue "Rodar as 5 execuções de cada versão e montar a tabela comparativa" "experimento,obrigatorio" "$M3" "$C" "$(cat <<'EOF'
Mínimo de 5 execuções por versão que conclui normalmente (V0, V2, V3).

Tabela com, por versão:
- [ ] tempo total de execução (média e desvio)
- [ ] tempo médio de espera pelos garfos
- [ ] maior tempo de espera observado
- [ ] número máximo de filósofos comendo simultaneamente
- [ ] ciclos concluídos por filósofo (todos devem dar 20)
- [ ] Salvar em `resultados/tabela.md` e no CSV
EOF
)"

issue "V2 – Investigar o impacto da sincronização excessiva" "v2,experimento,obrigatorio" "$M3" "$A" "$(cat <<'EOF'
O enunciado pede investigação **experimental** da V2 (destaque na página 2 do relatório).

- [ ] Comparar V2 com V0 e V3 usando os dados da tabela
- [ ] Mostrar que o máximo de filósofos comendo juntos na V2 é 1 (ou o limite da estratégia)
- [ ] Experimento extra sugerido: variar o tempo de comer (ex.: 5, 10, 20 ms) e mostrar como a diferença V2 × V3 cresce
- [ ] Texto curto com a conclusão, para usar nas perguntas 10 e 11
EOF
)"

issue "Gráficos comparativos" "experimento,relatorio,obrigatorio" "$M3" "$C" "$(cat <<'EOF'
- [ ] Barras: tempo total por versão (com desvio)
- [ ] Barras: espera média e espera máxima por versão
- [ ] Barras: máx. filósofos comendo simultaneamente
- [ ] (Opcional) gráfico do experimento extra da V2
- [ ] Salvar em `resultados/graficos/`
EOF
)"

# ===== 4. Relatório =====
issue "Relatório – Página 1: identificação, soluções, mecanismos e análise comparativa" "relatorio,obrigatorio" "$M4" "$C" "$(cat <<'EOF'
- [ ] Identificação do grupo (3 nomes)
- [ ] Descrição resumida de V0, V1, V2, V3 (textos vêm das issues de perguntas 1–3)
- [ ] Mecanismos de sincronização usados em cada versão
- [ ] Análise comparativa curta
EOF
)"

issue "Relatório – Página 2: tabela, gráficos e capturas de tela" "relatorio,obrigatorio" "$M4" "$C" "$(cat <<'EOF'
- [ ] Tabela de resultados (destaque para a V2)
- [ ] Gráficos comparativos
- [ ] Capturas de tela das execuções de V0, V2 e V3
- [ ] Evidência do deadlock da V1
EOF
)"

issue "Relatório – Perguntas 1 a 3 para V0 e V2" "relatorio,v0,v2,obrigatorio" "$M4" "$A" "$(cat <<'EOF'
1. Como funciona cada solução?
2. Quais mecanismos de sincronização são utilizados?
3. Como esses mecanismos são implementados em Go? (`sync.Mutex`, `sync.WaitGroup`, canais...)
EOF
)"

issue "Relatório – Perguntas 1 a 3 para V1 e V3" "relatorio,v1,v3,obrigatorio" "$M4" "$B" "$(cat <<'EOF'
1. Como funciona cada solução?
2. Quais mecanismos de sincronização são utilizados?
3. Como esses mecanismos são implementados em Go?
EOF
)"

issue "Relatório – Perguntas 4 a 6: deadlock e starvation" "relatorio,v1,v2,v3,obrigatorio" "$M4" "$B" "$(cat <<'EOF'
4. Por que ocorre deadlock na V1?
5. Como V2 e V3 evitam deadlock? (a parte da V2 é escrita por quem fez a V2 – Pessoa A)
6. Existe possibilidade de starvation? Justifique (para cada versão, com os ciclos por filósofo e a espera máxima)
EOF
)"

issue "Relatório – Perguntas 7 a 12: análise com as métricas" "relatorio,experimento,obrigatorio" "$M4" "$C" "$(cat <<'EOF'
Todas as respostas devem citar números da tabela, não só teoria.

7. Qual solução teve maior tempo total de execução?
8. Qual teve maior tempo de espera pelos recursos?
9. Quantos filósofos comeram simultaneamente em cada solução?
10. Qual foi o impacto da sincronização excessiva da V2?
11. A solução com mais sincronização foi a de melhor desempenho? Explique com as métricas.
12. Qual a relação entre grau de concorrência e tempo total?
EOF
)"

issue "Relatório – Conclusão, link do Git e exportação em PDF" "relatorio,obrigatorio" "$M4" "" "$(cat <<'EOF'
Tarefa de **todos**.
- [ ] Conclusão a partir da 3ª página, priorizando a análise dos resultados
- [ ] Link do repositório com os códigos completos e formatados
- [ ] Seção "Ambiente de execução"
- [ ] Conferir o `CHECKLIST.md` inteiro
- [ ] Exportar PDF e salvar em `relatorio/`
EOF
)"

# ===== 5. Apresentação =====
issue "Slides da apresentação" "apresentacao" "$M5" "$C" "$(cat <<'EOF'
- [ ] Problema e as 4 versões
- [ ] Deadlock da V1 (diagrama da espera circular)
- [ ] Tabela e gráficos
- [ ] Respostas-chave (perguntas 10–12)
- [ ] Dividir as falas entre os três
EOF
)"

issue "Demonstração ao vivo" "apresentacao,v1,v3" "$M5" "$B" "$(cat <<'EOF'
- [ ] Roteiro: V1 travando (modo forçado) → V2 lenta → V3 rápida
- [ ] Testar no notebook que será levado
- [ ] Ter prints e vídeo de reserva caso algo falhe
EOF
)"

issue "Ensaio geral – todos explicam todas as versões" "apresentacao" "$M5" "" "$(cat <<'EOF'
Os três precisam estar presentes, e o conteúdo cai na **Prova 1**.
- [ ] Cada integrante explica as 4 versões sem ler
- [ ] Cada integrante responde as 12 perguntas do relatório
- [ ] Ensaio cronometrado
EOF
)"

echo "Pronto: $(gh issue list --limit 100 --json number -q length) issues no repositório."
