install: clean
	go mod tidy
	cd cmd/varvoy && go install
	cd simdap && go install

clean:
	find . -name '_debug_bin*' | xargs rm -f
	find . -name 'varvoy_*' | xargs rm -rf
	find . -name '*.dot' | xargs rm -rf

lint:
	golangci-lint run