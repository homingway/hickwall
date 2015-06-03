#include <windows.h>
#include <stdio.h>
#include <pdh.h>
#include <pdhmsg.h>

#pragma comment(lib, "pdh.lib")

#ifndef PDH_FMT_NOCAP100
#define PDH_FMT_NOCAP100 ((DWORD) 0x00008000)
#endif

int getcpuload()
{
    static PDH_STATUS            status;
    static PDH_FMT_COUNTERVALUE  value;
    static HQUERY                query;
    static HCOUNTER              counter;
    static DWORD                 ret;
    static char                  runonce=1;
    char                         cput=0;

    if(runonce)
    {
        status = PdhOpenQuery(NULL, 0, &query);
        // printf("PdhOpenQuery: done: %x\n", query);
        if(status != ERROR_SUCCESS)
        {
            printf("PdhOpenQuery() ***Error: 0x%X\n",status);
            return 0;
        }

        // PdhAddCounter(query, TEXT("\\Processor(_Total)\\% Processor Time"),0,&counter); // A total of ALL CPU's in the system
        PdhAddCounter(query, TEXT("\\System\\Processes"),0,&counter); // A total of ALL CPU's in the system
        
        // PdhAddCounter(query, TEXT("\\Processes(_Total)\\Working Set"),0, &counter); // A total of ALL CPU's in the system
        
        //PdhAddCounter(query, TEXT("\\Processor(0)\\% Processor Time"),0,&counter);    // For systems with more than one CPU (Cpu0)
        //PdhAddCounter(query, TEXT("\\Processor(1)\\% Processor Time"),0,&counter);    // For systems with more than one CPU (Cpu1)
        runonce=0;
        PdhCollectQueryData(query); // No error checking here
        // return 0;
    }

    status = PdhCollectQueryData(query);
    if(status != ERROR_SUCCESS)
    {
        printf("PhdCollectQueryData() ***Error: 0x%X\n",status);
        if(status==PDH_INVALID_HANDLE) 
            printf("PDH_INVALID_HANDLE\n");
        else if(status==PDH_NO_DATA)
            printf("PDH_NO_DATA\n");
        else
            printf("Unknown error\n");
        return 0;
    }

    status = PdhGetFormattedCounterValue(counter, PDH_FMT_DOUBLE | PDH_FMT_NOCAP100 ,&ret, &value);
    if(status != ERROR_SUCCESS)
    {
        printf("PdhGetFormattedCounterValue() ***Error: 0x%X\n",status);
        return 0;
    }
    cput = value.doubleValue;

    return cput;
}

PDH_STATUS GetDoubleCounterValue(PDH_HCOUNTER counter,  double *value) {
    static PDH_FMT_COUNTERVALUE  cv;
    static PDH_STATUS            status;

    status = PdhGetFormattedCounterValue(counter, PDH_FMT_DOUBLE | PDH_FMT_NOCAP100,0, &cv);
    if(status == ERROR_SUCCESS) {
        *value = cv.doubleValue;    
    }
    return status;
}