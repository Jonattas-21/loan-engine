# Loan Engine

This project It's a credit simulator which allows users to consult loan conditions, based in some payments conditions.

## Features

- Loan simultor
- Loan conditions update
- Loan conditions Search

## Technical Componentes

- Authentication, prepared to use keycloak for authentication
- Caching stategy, to enhance performance while simulating loan or getting the conditions
- Simulator DB, stores the conditions to calculate the loan and the loan
- Logging storage, log separeted by level and prepared to quering in elastick search

## Instalation

1. Download the Go installer from the [official Go website](https://golang.org/dl/).
2. RInstall the linter, go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest.
3. Download `make` as part of the project.
4. Config the .env file, the example is in cmd/env_sample

### Configuring keycloak
1. Create a realm
2. Create a client with 'loan_app' name and client authorization property = off
3. Create a new user