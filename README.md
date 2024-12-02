# Projeto de Aprendizado de Golang

Este é um projeto desenvolvido para aprender a linguagem de programação Golang. O projeto utiliza a arquitetura Model, Repository, Controller e MongoDB como banco de dados.

## Estrutura do Projeto

- **Model**: Contém as definições das estruturas de dados utilizadas no projeto.
- **Repository**: Contém a lógica de acesso aos dados, incluindo operações de leitura e escrita no banco de dados MongoDB.
- **Controller**: Contém a lógica de controle, incluindo a manipulação das requisições HTTP e a interação com o repositório.

## Tecnologias Utilizadas

- **Golang**: Linguagem de programação utilizada para desenvolver o projeto.
- **MongoDB**: Banco de dados NoSQL utilizado para armazenar os dados.
- **godotenv**: Biblioteca utilizada para carregar variáveis de ambiente a partir de um arquivo `.env`.
- **bcrypt**: Biblioteca utilizada para criptografar senhas.
- **jwt-go**: Biblioteca utilizada para gerar e validar tokens JWT.

## Como Executar o Projeto

1. Clone o repositório:
    ```sh
    git clone <URL_DO_REPOSITORIO>
    ```

2. Navegue até o diretório do projeto:
    ```sh
    cd <NOME_DO_DIRETORIO>
    ```

3. Crie um arquivo `.env` com as seguintes variáveis:
    ```env
    MONGO_URI=<SUA_URI_DO_MONGODB>
    JWT_SECRET=<SEU_SEGREDO_JWT>
    ```

4. Execute o projeto:
    ```sh
    go run main.go
    ```

## Endpoints

- **/register**: Endpoint para registrar um novo usuário.
- **/login**: Endpoint para autenticar um usuário.
- **/friend**: Endpoint para gerenciar amigos de um usuário.