# Rinha de Back-end 2025

## Arquitetura

- 4 APIs recebem as requisições e adicionam os pagamentos em uma fila salva no Redis.
- 3 Workers de pagamento buscam pagamentos no Redis e os enviam aos Payment Processors.
- Enquanto as API e os Workers de pagamento trabalham, 1 Worker verifica de 5 em 5 segundos qual o melhor Payment Processor disponível.
- A abordagem de adicionar os pagamentos em uma fila garante que a aplicação tenha um baixo tempo de resposta enquanto os Workers verificam qual é o melhor Payment Processor disponível e enviam o pagamento para ele.
- A fila utiliza Mutex para garantir que os Workers não acessem o mesmo pagamento e gerem inconsistências.