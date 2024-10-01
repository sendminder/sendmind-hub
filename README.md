# sendmind-hub
this is sendmind core repo


## instal tools

1. brew install go-task/tap/go-task
2. go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## build step

1. cp .example.envrc .envrc
2. direnv allow
3. go mod download
4. task build


## push to docker hub

1. task push