all:
	nex rule.nex
	go tool yacc -o=rule.yacc.go rule.y
	go fmt 
	go build
test:
	go test -bench=Filter 
clean:
	-rm *.output *.yacc.go *.nn.go
