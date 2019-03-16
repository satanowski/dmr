build: clean
	go build -o cps

clean:
	@rm -f cps

cleanCache:
	@rm -f .usr.cache