go test . -gcflags -m -run "should_not_found_this_command" 2>&1 | grep -v "test" | grep -v "does not escape"
