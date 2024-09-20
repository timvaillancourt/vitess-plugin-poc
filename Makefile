all: cross_cell.so plugin-poc

cross_cell.so: durabler/cross_cell.go
	go build -buildmode=plugin ./durabler/cross_cell.go

plugin-poc: cross_cell.so durability.go main.go go.mod go.sum
	go build -o plugin-poc ./main.go ./durability.go
