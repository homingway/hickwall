go test . -gcflags -m -run Command 2>&1 | grep -v "test" | grep -v "does not escape"
