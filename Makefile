alibi: cmd/alibi/alibi.go
	go build -o $@ $^

run-alibi:
	go run cmd/alibi/alibi.go videodiary --missing
