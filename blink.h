#ifndef BLINK_DEFINE_H
#define BLINK_DEFINE_H

#include "stdio.h"
#include "mb.h"
#include "webview.h"
#include "util.h"

void mbInitialize();
void mbFinalize();
void mbSetMbDllPath(const char* dllPath);
void mbSetMbMainDllPath(const char* dllPath);

#endif