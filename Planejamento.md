# Planejamento da Iteração 1

## 1. Definir o foco da 1ª iteração

- CLI básico em Go (ex: comando `start`)
- Subir a aplicação Java (mesmo que fake)
- Implementar o `FakeSignatureService`
- Criar API `/sign` e `/validate` simples

---

## 2. Desenvolvendo um design (descrição de como vai funcionar)

### Java (primeira versão)
- Framework: Spring Boot  
- Endpoint: `POST /sign`  
- Recebe: string  
- Retorna: string fake (ex: `"signed-" + mensagem`)  

---

## 3. Planejamento da implementação

- Criar projeto Java (Spring Boot)  
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
