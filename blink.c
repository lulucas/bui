#include "blink.h"

typedef void (WKE_CALL_TYPE *FN_wkeInitializeEx)(const wkeSettings* settings);

__declspec(selectany) const wchar_t* s_wkeDllPath = L"ui.dll";
__declspec(selectany) HMODULE s_wkeMainDllHandle = NULL;

void wkeSetWkeDllHandle(const HMODULE mainDllHandle)
{
    s_wkeMainDllHandle = mainDllHandle;
}

void wkeSetWkeDllPath(const char* dllPath)
{
    size_t cSize = strlen(dllPath) + 1;
    wchar_t *wDllPath = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wDllPath, dllPath, cSize);

    s_wkeDllPath = wDllPath;
}

int wkeInitializeEx(const wkeSettings* settings)
{
    HMODULE hMod = s_wkeMainDllHandle;
    if (!hMod)
        hMod = LoadLibraryW(s_wkeDllPath);
    if (hMod) {
        FN_wkeInitializeEx wkeInitializeExFunc = (FN_wkeInitializeEx)GetProcAddress(hMod, "wkeInitializeEx");
        wkeInitializeExFunc(settings);

        WKE_FOR_EACH_DEFINE_FUNCTION(WKE_GET_PTR_ITERATOR0, WKE_GET_PTR_ITERATOR1, WKE_GET_PTR_ITERATOR2, WKE_GET_PTR_ITERATOR3, \
            WKE_GET_PTR_ITERATOR4, WKE_GET_PTR_ITERATOR5, WKE_GET_PTR_ITERATOR6, WKE_GET_PTR_ITERATOR11);
        return 1;
    }
    return 0;
}

int wkeInitialize()
{
    return wkeInitializeEx(((const wkeSettings*)0));
}