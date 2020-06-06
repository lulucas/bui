#include "util.h"
#include <stdio.h>

const wchar_t* utf8ToWchar(const char* utf8Str) {
    // UTF8 TO Wchar
    size_t n = MultiByteToWideChar(CP_UTF8, 0, utf8Str, -1, NULL, 0);
    if (0 == n) return L"";
    wchar_t* wbuf = (wchar_t *)malloc(sizeof(wchar_t)*(n+1));
    MultiByteToWideChar(CP_UTF8, 0, utf8Str, -1, wbuf, n);
    return wbuf;
}
