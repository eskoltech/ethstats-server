b:
	go build -o ./build/bin/ethstats

start:
	./build/bin/ethstats --secret 123456789
