#include "resource.h"
#include <signal.h>
#include <windows.h>
#include <tchar.h>
#include <stdio.h>

#include <pdh.h>
#include <pdhmsg.h>

#ifdef _MSC_VER
#pragma comment(lib, "pdh.lib")
#endif

#ifndef PDH_FMT_NOCAP100
#define PDH_FMT_NOCAP100 ((DWORD) 0x00008000)
#endif


#ifdef __cplusplus
extern "C" {
#endif


int cpuusage();

int getcpuload();


#ifdef __cplusplus
}
#endif