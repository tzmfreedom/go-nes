FILE := sample1.nes

.PHONY: run
run:
	go run . $(FILE)
