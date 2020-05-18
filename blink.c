#include "blink.h"

// mb

typedef void (MB_CALL_TYPE *FN_mbInit)(const mbSettings* settings);

__declspec(selectany) const wchar_t* kMbDllPath = L"mb.dll";
__declspec(selectany) const wchar_t* kMbMainDllPath = L"node.dll";

void mbSetMbDllPath(const char* dllPath)
{
    kMbDllPath = utf8ToUtf16(dllPath);
}

void mbSetMbMainDllPath(const char* dllPath)
{
    kMbMainDllPath = utf8ToUtf16(dllPath);
}

void mbInit(const mbSettings* settings)
{
    LoadLibraryW(kMbMainDllPath);
    HMODULE hMod = LoadLibraryW(kMbDllPath);

    FN_mbInit mbInitExFunc = (FN_mbInit)GetProcAddress(hMod, "mbInit");
    mbInitExFunc(settings);

    MB_FOR_EACH_DEFINE_FUNCTION(MB_GET_PTR_ITERATOR0, MB_GET_PTR_ITERATOR1, MB_GET_PTR_ITERATOR2, MB_GET_PTR_ITERATOR3, \
        MB_GET_PTR_ITERATOR4, MB_GET_PTR_ITERATOR5, MB_GET_PTR_ITERATOR6, MB_GET_PTR_ITERATOR7, MB_GET_PTR_ITERATOR8, MB_GET_PTR_ITERATOR9, MB_GET_PTR_ITERATOR10, MB_GET_PTR_ITERATOR11);

    return;
}

void mbInitialize()
{
    mbSettings settings;
    memset(&settings, 0, sizeof(settings));
    mbInit(&settings);
}

void mbFinalize() {
    mbUninit();
}