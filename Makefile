gen:
	@go build -o bin/gen data_gen/*.go
	@./bin/gen

rec:
	@go build -o bin/recv data_recv/*.go
	@./bin/recv