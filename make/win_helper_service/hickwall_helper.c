#include <Windows.h>
#include <stdio.h>

#define SLEEP_TIME 1000
#define SVC_NAME "hickwallhelper"

SERVICE_STATUS ServiceStatus;
SERVICE_STATUS_HANDLE hStatus;

void  ServiceMain(int argc, char** argv);
void  CtrlHandler(DWORD request);
int InitService();
void start_hickwall_if_not_running();

// Service initialization
int InitService() 
{ 
    int result;
    return(result); 
}

// Control Handler
void  CtrlHandler(DWORD request)
{
    switch (request)
    {
    case SERVICE_CONTROL_STOP:
        ServiceStatus.dwWin32ExitCode = 0; 
        ServiceStatus.dwCurrentState = SERVICE_STOPPED; 
        SetServiceStatus (hStatus, &ServiceStatus);
        return;
    case SERVICE_CONTROL_SHUTDOWN:
        ServiceStatus.dwWin32ExitCode = 0; 
        ServiceStatus.dwCurrentState = SERVICE_STOPPED; 
        SetServiceStatus (hStatus, &ServiceStatus);
        return;
    default:
        break;
    }
    // Report current status
    SetServiceStatus (hStatus, &ServiceStatus);
    return;
}

void  ServiceMain(int argc, char** argv)
{
    int error;
    ServiceStatus.dwServiceType = 
        SERVICE_WIN32;
    ServiceStatus.dwCurrentState = 
        SERVICE_START_PENDING;
    ServiceStatus.dwControlsAccepted = 
        SERVICE_ACCEPT_SHUTDOWN | 
        SERVICE_ACCEPT_STOP;
    ServiceStatus.dwWin32ExitCode = 0;
    ServiceStatus.dwServiceSpecificExitCode = 0;
    ServiceStatus.dwCheckPoint = 0;
    ServiceStatus.dwWaitHint = 0;

    hStatus = RegisterServiceCtrlHandler(SVC_NAME, (LPHANDLER_FUNCTION)CtrlHandler);
    if (hStatus == (SERVICE_STATUS_HANDLE)0)
    {
        return;
    }
    // Initialize Service 
    error = InitService();
    if (error) 
    {
        // Initialization failed
        ServiceStatus.dwCurrentState = 
            SERVICE_STOPPED; 
        ServiceStatus.dwWin32ExitCode = -1; 
        SetServiceStatus(hStatus, &ServiceStatus); 
        return; 
    } 

    ServiceStatus.dwCurrentState = 
        SERVICE_RUNNING;
    SetServiceStatus (hStatus, &ServiceStatus);
    
    MEMORYSTATUS memstatus;
    while (ServiceStatus.dwCurrentState == SERVICE_RUNNING)
    {
        Sleep(SLEEP_TIME);
        start_hickwall_if_not_running();
    }
}

void start_hickwall_if_not_running()
{
    SC_HANDLE schSCManager;
    SC_HANDLE schService;
    SERVICE_STATUS  lpServiceStatus;

    schSCManager = OpenSCManager(NULL, NULL, SC_MANAGER_ALL_ACCESS);
    if (NULL == schSCManager){
        return;
    }

    // Get a handle to the service.
    schService = OpenService(schSCManager, "hickwall", SC_MANAGER_ALL_ACCESS);
     if (schService == NULL)
    { 
        CloseServiceHandle(schSCManager);
        return;
    }    

    if (!QueryServiceStatus(schService, &lpServiceStatus)){
        CloseServiceHandle(schService); 
        CloseServiceHandle(schSCManager);
        return;
    }

    if (lpServiceStatus.dwCurrentState == SERVICE_STOPPED)
    {
        StartService(schService, 0, NULL);
    }

    CloseServiceHandle(schService); 
    CloseServiceHandle(schSCManager);
}


int main()
{
    SERVICE_TABLE_ENTRY ServiceTable[2];
    ServiceTable[0].lpServiceName=SVC_NAME;
    ServiceTable[0].lpServiceProc=(LPSERVICE_MAIN_FUNCTION)ServiceMain;
    ServiceTable[1].lpServiceName=NULL;
    ServiceTable[1].lpServiceProc=NULL;
    StartServiceCtrlDispatcher(ServiceTable);
}