# Checklist de entrega – tudo o que o enunciado cobra

Cada linha aponta para a issue que cobre o item. Antes de entregar, **todas** precisam estar marcadas.

## Implementações (Seção 2)
- [x] V0 sequencial, baseada na versão do Moodle → *V0 – Versão sequencial*
- [x] V1 com deadlock, executada e demonstrada → *V1 – Versão concorrente sujeita a deadlock*
- [x] V1: recursos que cada filósofo possui → *V1 – Evidência e análise*
- [x] V1: recurso pelo qual cada um espera → *V1 – Evidência e análise*
- [x] V1: como ocorre a espera circular → *V1 – Evidência e análise*
- [x] V1: as 4 condições (exclusão mútua, segura-e-espera, sem preempção, espera circular) → *V1 – Evidência e análise*
- [x] V2 correta, sem deadlock, com concorrência limitada → *V2 – Sem deadlock com sincronização excessiva*
- [x] V2: impacto investigado **experimentalmente** com medição de tempo → *V2 – Investigar o impacto*
- [x] V3 sem deadlock e com mais concorrência, estratégia justificada → *V3 – Versão concorrente aprimorada*
- [x] V3: por que evita deadlock → *V3*
- [x] V3: discussão sobre starvation → *V3*
- [x] Mecanismos sem equivalente nativo em Go adaptados/refatorados e explicados → *V2 / V3*

## Experimento (Seção 3)
- [x] 5 filósofos, 20 ciclos cada: pensar → adquirir → comer → liberar → *Fixar parâmetros*
- [x] Mesmos tempos de pensar/comer em todas as versões → *Fixar parâmetros*
- [x] ≥ 5 execuções de cada versão que termina → *Rodar as 5 execuções*
- [x] Tempo total de execução → *Tabela comparativa*
- [x] Tempo médio de espera pelos garfos → *Tabela comparativa*
- [x] Maior tempo de espera → *Tabela comparativa*
- [x] Máx. de filósofos comendo simultaneamente → *Tabela comparativa*
- [x] Ciclos concluídos por filósofo → *Tabela comparativa*
- [x] Resultados em tabela comparativa → *Tabela comparativa*
- [x] V1 analisada separadamente com evidência do bloqueio → *V1 – Evidência e análise*

## Relatório em PDF (Seção 4)
- [x] **Pág. 1:** identificação do grupo, descrição das soluções, mecanismos, análise comparativa
- [x] **Pág. 2:** tabela (destaque V2), gráficos, prints das execuções, evidência do deadlock da V1
- [x] **Pág. 3+:** conclusão baseada nas métricas
- [x] Link do Git com códigos completos e formatados (`gofmt`)
- [x] Perguntas 1–12 todas respondidas, com números:
  - [x] 1. Como funciona cada solução
  - [x] 2. Mecanismos de sincronização
  - [x] 3. Como são implementados em Go
  - [x] 4. Por que ocorre deadlock na V1
  - [x] 5. Como V2 e V3 evitam deadlock
  - [x] 6. Starvation – justificar
  - [x] 7. Maior tempo total
  - [x] 8. Maior tempo de espera
  - [x] 9. Quantos comeram simultaneamente em cada solução
  - [x] 10. Impacto da sincronização excessiva da V2
  - [x] 11. Mais sincronização = melhor desempenho? (com métricas)
  - [x] 12. Relação grau de concorrência × tempo total
- [x] Ambiente: linguagem e versão, SO, processador, nº de cores, comandos de compilação/execução
- [x] Máquina multicore

