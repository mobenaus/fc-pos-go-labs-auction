# fc-pos-go-labs-auction

## Descrição do Desafio
```
Objetivo: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: [clique para acessar o repositório](https://github.com/devfullcycle/labs-auction-goexpert).

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

Você deverá desenvolver:

Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
Um teste para validar se o fechamento está acontecendo de forma automatizada;

Dicas:

Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;
 
Entrega:

O código-fonte completo da implementação.
Documentação explicando como rodar o projeto em ambiente dev.
Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.
```

## Execução da solução

```
docker compose up
```
O docker compose é formado pelos seguintes serviços:
- mongodb: banco de dados
- mongo-express: serviço web para acessar o mongodb
- auction: serviço com a implementação do serviço de leilão, com as alterações para fechar o leilão automaticamente (configurado em 20s)
- auction_test: container que inicia e executa os testes do serviço auction, os testes executados são:
  - chamar endpoint para criar um usuário
  - chamar endpoint para criar um leilão (auction)
  - chamar endpoint para buscar o ultimo leilão aberto e recuperar o Id
  - chamar 30 vezes o endpoint de lance (bid) 1 vez por segundo, nem todos serão registrados devido ao fechamento automatico do leilão
  - chamar endpoint de consulta de auction e verificar o estado do leilao (2-fechado)
  - chamar endpoing de consulta de ganhador do leição (auction/winner)
  
## Alterações efetuadas na solução

1. Adicionado mongo-express no docker compose
2. Implementado o metodo completeAuction e getAuctionInterval em create_auction.go
3. Otimização de Dockerfile
4. Implementação dos testes em auction_test
5. Ajustes no logger, AuctionStatus e outros
6. Implementação de CreateUser
7. Adicionado auction_test no docker compose