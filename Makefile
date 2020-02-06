fetch-episodes:
	node -r dotenv/config ./scripts/fetchEpisodes.js

gofmt:
	gofmt -w internal
	gofmt -w pkg

golint:
	golint internal
	golint pkg

tryout:
	go run ./pkg/main.go

alibi: cmd/alibi/alibi.go
	go build -o $@ $^

run-alibi:
	go run cmd/alibi/alibi.go videodiary --missing
