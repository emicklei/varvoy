install:
	cd cmd/varvoy && go install
	cd simdap && go install

clean:
	find . -name '_debug_bin*' | xargs rm -f