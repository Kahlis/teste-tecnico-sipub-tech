## Visão Geral
> O sistema visa oferecer uma API REST para gerenciamento de filmes, utilizando arquitetura hexagonal e comunicação entre microserviços via gRPC.

---

## Requisitos Funcionais (RF)

| ID    | Nome         | Descrição                                         | Prioridade                          | Critério de Aceitação                                                    |
| ----- | ------------ | ------------------------------------------------- | ----------------------------------- | ------------------------------------------------------------------------ |
| RF-01 | **[Read]**   | Deve ser possível ler um único filme              | Alta                                | GET `/movies/{id}` deve retornar o filme recuperado do banco pelo `{id}` |
| RF-02 | **[Read N]** | Deve ser possível listar todos os filmes          | Alta                                | GET `/movies` deve retornar uma lista dos filmes disponíveis no banco    |
| RF-03 | **[Delete]** | Deve ser possível atualizar o cadastro dos filmes | Alta                                | ...                                                                      |
| RF-04 | **[Update]** | Deve ser possível atualizar o cadastro dos filmes | Baixa<br>**(Não foi especificado)** | PATCH `/movies/{id}` deve atualizar o cadastro do filme de ID `{id}`     |


---

## Requisitos Não Funcionais (RNF)

| ID     | Categoria                      | Descrição                                                                                        | Critério de Avaliação                                                                                           |
| ------ | ------------------------------ | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------- |
| RNF-01 | **Arquitetura**                | Deve usar microsserviços com containeres para **API Gateway**, **Movies** e **Banco de dados**   | Cada serviço rodará em seu respectivo container                                                                 |
| RNF-02 | **Comunicação**                | A comunicação entre **API** e **Movies** deve ser feita via **gRPC Protocol Buffers (Protobuf)** | Verificar que as chamadas entre os serviços ocorrem via gRPC utilizando arquivos `.proto` devidamente definidos |
| RNF-03 | **Usabilidade / Documentação** | A documentação da **API** deve ser feita utilizando **Swagger**                                  | Cada rota deve ter documentação adequada via swagger                                                            |

---
## Restrições do Projeto
- [ ] O projeto deve ser desenvolvido utilizando **Go (Golang)**;
- [ ] O banco de dados deve ser **MongoDB** ou **DynamoDB** (emulado via LocalStack, se aplicável);
- [ ] O código-fonte deve ser versionado e disponibilizado em um **repositório público no GitHub**;
- [ ] O link do repositório deve ser enviado por e-mail para **[rh@facilita.tech](mailto:rh@facilita.tech)** até **01/11/2025**;
- [ ] O título do e-mail deve ser exatamente (em letras maiúsculas):  `TESTE TÉCNICO CONCLUIDO - SIPUB TECH`
