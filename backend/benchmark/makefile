.PHONY: clean bench

clean:
	rm ../app/*.txt
	rm ./*.txt

bench:
	make clean
	@echo "Starting docker stats logging..."
	@docker stats >> "stats.txt" & echo $$! > stats.pid
	@echo "Running Go application..."
	@go run ./
	@sleep 10
	@echo "Stopping docker stats logging..."
	@kill $$(cat stats.pid) && rm stats.pid
	@echo "Benchmark completed."
