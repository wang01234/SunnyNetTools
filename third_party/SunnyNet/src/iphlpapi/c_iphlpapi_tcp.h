#include <windows.h>
#include <iphlpapi.h>
#include <stdio.h>

typedef DWORD(WINAPI* GETEXTENDEDTABLE)(PVOID, PDWORD, BOOL, ULONG, TCP_TABLE_CLASS, ULONG);
typedef DWORD(WINAPI* SETTCPENTRY)(PMIB_TCPROW);

void closeTcpConnectionInit();
void closeTcpConnectionByPid(DWORD pid, DWORD ulAf);
int getTcpInfoPID(char* Addr, int SunnyProt);
