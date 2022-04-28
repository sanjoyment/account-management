cert:
	cd cert; chmod +x gen.sh; ./gen.sh; cd ..

run:
	go run main.go

.PHONY: gen clean server client test cert