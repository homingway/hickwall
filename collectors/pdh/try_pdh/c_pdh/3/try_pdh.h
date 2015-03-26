#include <windows.h>
#include <stdio.h>
#include <pdh.h>
#include <pdhmsg.h>

#pragma comment(lib, "pdh.lib")

#ifndef PDH_FMT_NOCAP100
#define PDH_FMT_NOCAP100 ((DWORD) 0x00008000)
#endif

int getcpuload();
