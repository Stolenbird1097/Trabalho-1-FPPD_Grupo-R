# Checklist de entrega – tudo o que o enunciado cobra

Cada linha aponta para a issue que cobre o item. Antes de entregar, **todas** precisam estar marcadas.

## Implementações (Seção 2)
- [ ] V0 sequencial, baseada na versão do Moodle → *V0 – Versão sequencial*
- [ ] V1 com deadlock, executada e demonstrada → *V1 – Versão concorrente sujeita a deadlock*
- [ ] V1: recursos que cada filósofo possui → *V1 – Evidência e análise*
- [ ] V1: recurso pelo qual cada um espera → *V1 – Evidência e análise*
- [ ] V1: como ocorre a espera circular → *V1 – Evidência e análise*
- [ ] V1: as 4 condições (exclusão mútua, segura-e-espera, sem preempção, espera circular) → *V1 – Evidência e análise*
- [ ] V2 correta, sem deadlock, com concorrência limitada → *V2 – Sem deadlock com sincronização excessiva*
- [ ] V2: impacto investigado **experimentalmente** com medição de tempo → *V2 – Investigar o impacto*
- [ ] V3 sem deadlock e com mais concorrência, estratégia justificada → *V3 – Versão concorrente aprimorada*
- [ ] V3: por que evita deadlock → *V3*
- [ ] V3: discussão sobre starvation → *V3*
- [ ] Mecanismos sem equivalente nativo em Go adaptados/refatorados e explicados → *V2 / V3*

## Experimento (Seção 3)
- [ ] 5 filósofos, 20 ciclos cada: pensar → adquirir → comer → liberar → *Fixar parâmetros*
- [ ] Mesmos tempos de pensar/comer em todas as versões → *Fixar parâmetros*
- [ ] ≥ 5 execuções de cada versão que termina → *Rodar as 5 execuções*
- [ ] Tempo total de execução → *Tabela comparativa*
- [ ] Tempo médio de espera pelos garfos → *Tabela comparativa*
- [ ] Maior tempo de espera → *Tabela comparativa*
- [ ] Máx. de filósofos comendo simultaneamente → *Tabela comparativa*
- [ ] Ciclos concluídos por filósofo → *Tabela comparativa*
- [ ] Resultados em tabela comparativa → *Tabela comparativa*
- [ ] V1 analisada separadamente com evidência do bloqueio → *V1 – Evidência e análise*

## Relatório em PDF (Seção 4)
- [ ] **Pág. 1:** identificação do grupo, descrição das soluções, mecanismos, análise comparativa
- [ ] **Pág. 2:** tabela (destaque V2), gráficos, prints das execuções, evidência do deadlock da V1
- [ ] **Pág. 3+:** conclusão baseada nas métricas
- [ ] Link do Git com códigos completos e formatados (`gofmt`)
- [ ] Perguntas 1–12 todas respondidas, com números:
  - [ ] 1. Como funciona cada solução
  - [ ] 2. Mecanismos de sincronização
  - [ ] 3. Como são implementados em Go
  - [ ] 4. Por que ocorre deadlock na V1
  - [ ] 5. Como V2 e V3 evitam deadlock
  - [ ] 6. Starvation – justificar
  - [ ] 7. Maior tempo total
  - [ ] 8. Maior tempo de espera
  - [ ] 9. Quantos comeram simultaneamente em cada solução
  - [ ] 10. Impacto da sincronização excessiva da V2
  - [ ] 11. Mais sincronização = melhor desempenho? (com métricas)
  - [ ] 12. Relação grau de concorrência × tempo total
- [ ] Ambiente: linguagem e versão, SO, processador, nº de cores, comandos de compilação/execução
- [ ] Máquina multicore

## Apresentação (50% da nota)
- [ ] Os três presentes
- [ ] Demo: V1 travando, V2 e V3 rodando
- [ ] Todos sabem explicar tudo (conteúdo da Prova 1)
