# go-docker

Containerizing a Go API

## Usage

`make cert`

`make api`

```
├── .dockerignore
├── Dockerfile
├── README.md
├── cmd
|  └── web
|     ├── handlers.go
|     ├── helpers.go
|     ├── main.go
|     ├── middleware.go
|     └── routes.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── makefile
└── tls
   ├── cert.pem
   └── key.pem
```
