# Loan Engine

This project It's a credit simulator which allows users to consult loan conditions, based in some payments conditions.

## Features

- Loan simultor
- Interest loan conditions update by age group
- Get all loan age group interest
- Send email for each loan simulation

### Activity Diagram
Bellow folow two use cases that illustrate what it's possible to operate in the system.

![alt text](https://raw.githubusercontent.com/Jonattas-21/loan-engine/refs/heads/main/docs_internal/activity_diagram.png?token=GHSAT0AAAAAAC5MOQWOEDLPKWOHAJG6SUNGZ44JR6Q "USe Case")

### Container Diagram
The system has been splitted in five distinguiguished services to manage each specific responsability.

![alt text](https://raw.githubusercontent.com/Jonattas-21/loan-engine/refs/heads/main/docs_internal/topology.png?token=GHSAT0AAAAAAC5MOQWOXL47BLS4JICI7EX6Z44JS5Q "Container Diagram")

## Technical Componentes

- Authentication, ready to use keycloak for authentication using oAuth 2
- Caching stategy, to enhance performance while simulating loan or getting the conditions
- Simulator DB, stores the conditions to calculate the loan and the loan simulated
- Logging storage, log separeted by level and prepared to quering in elastick search
- Swagger for api documentation

## Instalation

1. Download the Go installer from the [official Go website](https://golang.org/dl/).
2. Install the linter, go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest.
3. Download `make` as part of the project.
4. Config the .env file, the example is in cmd/env_sample
5. >make run

### Configuring keycloak
1. Create a realm and check the url: KEYCLOAK_HOST="http://localhost:8080/realms/{realmName}" in .env
2. Create a client with 'loan_app' name and client authorization property = off in advance tab
3. Create a new user and password
