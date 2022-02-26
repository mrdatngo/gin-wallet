## Golang Gin Framework Fundamental

Golang using gin framework - Wallet Management

* ### Usage
  - This project is build with Wallet Management
  - - Step to use:
  1. Install mysql & create schema ```wallet```, after that run project to create tables
  2. Import gin-api.postman_collection.json to have apis for postman test
  3. Run project: ``` go run main.go ```
  4. Register a user: localhost:3001/api/v1/register, should register two (1 admin, 1 merchant user). To make user become admin, go to mysql db and change RoleID to 1
  5. With dev environment, register api will return token in data response, take it to activation account: ```localhost:3001/api/v1/activation/:token```
  6. To set up user have Admin permission, go to db and edit RoleID of user to 1
  7. Login admin user (roleID = 1) to get token authen: ```http://localhost:3001/api/v1/login```
  8. By default, user have no Wallet, create wallet with default balance 0 (admin role): ```localhost:3001/api/v1/wallet```
  9. To view all wallet of user (users themselves or admin role), use api: ```localhost:3001/api/v1/wallets```
  10. Can delete wallet at (only admin role): ```localhost:3001/api/v1/wallet/:wallet_id```, beware that wallet not actually delete, it just de-activate (we don't delete data on real-case)
  11. Deposit to wallet (admin role): ```localhost:3001/api/v1/deposit```
  12. View deposits (admin role or users themselves): ```localhost:3001/api/v1/deposits```

### Feature

- [x] Containerize Application Using Docker
- [x] Protected Route Using JWT
- [x] Integerate ORM Database Using Gorm
- [x] API Documentation Using Swagger
- [x] Validation Request Using Go Playground Validator
- [x] Integerate Unit Testing
- [x] And More

## Command

- ### Application Lifecycle

  - Install node modules

  ```sh
  $ go get . || go mod || make goinstall
  ```

  - Build application

  ```sh
  $ go build -o main || make goprod
  ```

  - Start application in development

  ```sh
  $ go run main.go | make godev
  ```

  - Test application

  ```sh
  $ go test main.go main_test.go || make gotest
  ```