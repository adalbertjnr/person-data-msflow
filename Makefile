gen:
	@go build -o bin/gen data_gen/*.go
	@./bin/gen

rec:
	@go build -o bin/recv data_recv/*.go
	@./bin/recv

prc:
	@go build -o bin/prc data_proc/*.go
	@./bin/prc

agg:
	@go build -o bin/agg aggregator/*.go
	@./bin/agg