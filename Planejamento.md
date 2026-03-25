# Antes da Iteração 1
## O que estamos construindo?
Temos essencialmente três coisas:

- Dois CLIs (Command Line Interfaces)
- Uma aplicação Java (assinador/simulador)

O CLI não faz o trabalho principal — apenas chama as funções de outras aplicações (assinador, hubsaude).
As responsabilidades do CLI incluem gerenciar o ciclo de vida das aplicações e facilitar seu uso.

Enquanto isso, o assinador.jar:
- faz assinatura (simula)
- possivelmente roda como servidor HTTP

# Planejamento da Iteração 1

## 1. Definir o foco da 1ª iteração

- CLI básico em Go (comando `start`)
- Subir a aplicação Java (fake)
- Implementar o `FakeSignatureService`
- Criar API `/sign` e `/validate` simples

---

## 2. Desenvolvendo um design (descrição de como vai funcionar)

### Java (primeira versão)
- Framework: Spring Boot  
- Endpoint: `POST /sign`  
- Recebe: string  
- Retorna: string fake 

---

## 3. Planejamento da implementação

- Criar projeto Java 
- Criar controller `/sign`  
- Criar `SignatureService`  
- Criar `FakeSignatureService`  
- Testar endpoint com Postman  

---

## 4. Definir testes

- Testar se `/sign` retorna algo  
- Testar se não quebra com entrada vazia  

---

## 5. Definir ambiente

- Código no GitHub  
- Branches: `main`, `feature/...`  
- Linguagens: Go (CLI) e Java (backend)  

---

## 6. Descrever o ciclo da iteração

### Iteração 1:
- Design: API de assinatura fake  
- Implementação: endpoint `/sign`  
- Testes: validar resposta  
- Revisão: código do grupo  
- Refatoração: ajustes simples  
