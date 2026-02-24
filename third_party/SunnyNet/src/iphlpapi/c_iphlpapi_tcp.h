#include <windows.h>
#include <iphlpapi.h>
#include <stdio.h>

// 使用 Windows SDK 中已定义的结构体，不重复定义
// MIB_TCPROW2 和 MIB_TCPTABLE2 已由 iphlpapi.h 定义

// 定义 GETEXTENDEDTABLE 函数指针类型
typedef DWORD(WINAPI* GETEXTENDEDTABLE)(PVOID, PDWORD, BOOL, ULONG, TCP_TABLE_CLASS, ULONG);

// 定义 SETTCPENTRY 函数指针类型
typedef DWORD(WINAPI* SETTCPENTRY)(PMIB_TCPROW);

// 定义 GetTcpTable2 函数指针类型 (使用 SDK 中的定义)
typedef DWORD (WINAPI * GetTcpTable2)(PMIB_TCPTABLE2 TcpTable, PULONG SizePointer, BOOL Order);

// 关闭 TCP 连接初始化
void closeTcpConnectionInit();

// 根据 PID 关闭 TCP 连接
void closeTcpConnectionByPid(DWORD pid, DWORD ulAf);

// 获取指定 TCP 地址和端口的 PID
int getTcpInfoPID(char* Addr, int SunnyProt);
