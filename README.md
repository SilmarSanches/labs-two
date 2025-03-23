# labs-two

Projeto de conclusão de pós-graduação

## Indice
1. [Build-Dev](#build)
2. [DockerCompose](#dockercompose)
3. [Testes](#testes)
4. [Swagger](#swagger)

## Build-Dev

O serviceA é acessível na porta 8080. Já o swerviceB é acessivel na porta 8081.

- ServiceA: entrar na pasta serviceA e rodar o comando
```bash
go run ./cmd
```
- ServiceB: entrar na pasta serviceB e rodar o comando
```bash
go run ./cmd
```
## DockerCompose

O docker compose está na raiz do projeto, então rodar o comando abaixo para a execução:

```bash
docker compose up -d
````

## Testes

Tanto no projeto serviceA quando no projeto serviceB existe a pasta api contendo os testes.

Alem disso existem alguns testes unitários que podem ser rodados através do comando:

- ServiceA:
```bash
cd serviceA
go test ./...
````

- ServiceB:
```bash
cd serviceB
go test ./...
````

## Swagger
- o swagger pode ser acessado localment na porta 8080 ou 8081 respectivamente:

```bash
http://localhost:8080/swagger/index.html

ou

http://localhost:8081/swagger/index.html
```