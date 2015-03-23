// +build windows

package main

// import (
// 	"fmt"
// 	"syscall"
// 	"unsafe"
// )

/*
#include <windows.h>
#include <stdio.h>
#include <_mingw_mac.h>
#include <pdh.h>
#include <pdhmsg.h>

#pragma comment(lib, "pdh.lib")


void Func1(){
    HQUERY hQuery = NULL;
    HLOG hLog = NULL;
    PDH_STATUS pdhStatus;
    DWORD dwLogType = PDH_LOG_TYPE_CSV;
    HCOUNTER hCounter;
    DWORD dwCount;

    printf("%s", "this is a test\n");

    pdhStatus = PdhOpenQuery(NULL, 0, &hQuery);
}

*/
import "C"

func main() {
	C.Func1()
}
