package main

import (
	"fmt"
	"os"
)

/*
#include<windows.h>
#include<stdio.h>
#include<tchar.h>
#include<psapi.h>

#pragma comment(lib, "pdh.lib")

static ULARGE_INTEGER lastCPU, lastSysCPU, lastUserCPU;
static int numProcessors;
static HANDLE self;


void init(){
    SYSTEM_INFO sysInfo;
    FILETIME ftime, fsys, fuser;


    GetSystemInfo(&sysInfo);
    numProcessors = sysInfo.dwNumberOfProcessors;


    GetSystemTimeAsFileTime(&ftime);
    memcpy(&lastCPU, &ftime, sizeof(FILETIME));


    self = GetCurrentProcess();
    GetProcessTimes(self, &ftime, &ftime, &fsys, &fuser);
    memcpy(&lastSysCPU, &fsys, sizeof(FILETIME));
    memcpy(&lastUserCPU, &fuser, sizeof(FILETIME));
}

DWORDLONG GetMemTotal(){
    MEMORYSTATUSEX statex;
    statex.dwLength = sizeof(statex);
    GlobalMemoryStatusEx(&statex);
    return statex.ullTotalPhys;
}

DWORDLONG GetCurrentValue(){
    PROCESS_MEMORY_COUNTERS pmc;
    GetProcessMemoryInfo(GetCurrentProcess(), &pmc, sizeof(pmc));
    //SIZE_T virtualMemUsedByMe = pmc.PrivateUsage;
    return pmc.WorkingSetSize;
}

*/
import "C"

func main() {
	h, _ := os.Hostname()
	fmt.Println(h)
	fmt.Println("", C.GetMemTotal()/1024/1024)
	fmt.Println(C.GetCurrentValue())
}
