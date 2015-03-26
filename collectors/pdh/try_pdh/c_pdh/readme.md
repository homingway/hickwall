build notes:
==========


####First Method:

    go build try_pdh.go
    
use this method, we cann directly `#include "try_pdh.c"`  **[see example in folder 1]**

####Second Method:

    go build -o try_pdh.exe

this method is a little bit tricky.

* we have to change `try_pdh.c` into `try_pdh.h` and include .h file. otherwise, multiple defination error will be raised. because go will try to build `try_pdh.c` into `try_pdh.o` first, and then build `try_pdh.go` with link to `try_pdh.o`.  all definations will be duplicated this way. **[see example in folder 2]**

* or we can split `try_pdh.c` into 2 file: header and c content. and include header file within both c conent file and go file. **[see example in folder 3]**