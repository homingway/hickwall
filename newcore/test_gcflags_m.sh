go test . -run "should_not_found_this_command" -gcflags -m 2>&1 | grep -v "test" | grep -v "does not escape"
go test . -run "should_not_found_this_command" -gcflags -m 2>&1 | grep -v "test" | grep -v "does not escape" | wc -l
