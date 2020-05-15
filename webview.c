#include "webview.h"

int goGetBuiPort(wkeWebView window);
jsValue WKE_CALL_TYPE onBuiPort(jsExecState es, void *param)
{
    int port = goGetBuiPort(jsGetWebView(es));
    return jsInt(port);
}
void bindPort() {
    wkeJsBindFunction("BUI_PORT", onBuiPort, NULL, 2);
}

wkeWebView createWebWindow(int width, int height, bool transparent) {
    wkeWebView window = wkeCreateWebWindow(transparent ? WKE_WINDOW_TYPE_TRANSPARENT : WKE_WINDOW_TYPE_POPUP, NULL, 0, 0, width, height);
    return window;
}

HWND getWindowHandle(wkeWebView window)
{
    return wkeGetWindowHandle(window);
}

void loadURL(wkeWebView window, char *url)
{
    wkeLoadURL(window, url);
    free(url);
}

void reloadURL(wkeWebView window)
{
    wkeReload(window);
}

void setWindowTitle(wkeWebView window, char *title)
{
    wkeSetWindowTitle(window, title);
    free(title);
}

void destroyWindow(wkeWebView window)
{
    wkeDestroyWebWindow(window);
}

void showWindow(wkeWebView window, bool show) {
    wkeShowWindow(window, show);
}

void moveToCenter(wkeWebView window) {
    wkeMoveToCenter(window);
}

void setLocalStorageFullPath(wkeWebView webView, const char* path) {
    size_t cSize = strlen(path) + 1;
    wchar_t *wPath = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wPath, path, cSize);
    wkeSetLocalStorageFullPath(webView, wPath);
}

void setCookieJarFullPath(wkeWebView webView, const char* path) {
    size_t cSize = strlen(path) + 1;
    wchar_t *wPath = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wPath, path, cSize);
    wkeSetCookieJarFullPath(webView, wPath);
}

void showDevtools(wkeWebView webView, const char* path) {
    size_t cSize = strlen(path) + 1;
    wchar_t *wPath = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wPath, path, cSize);
    wkeShowDevtools(webView, wPath, NULL, NULL);
}

// ----------- Callback -------------

void goOnDocumentReady(wkeWebView window, void *param);
void onDocumentReady(wkeWebView window, void* param) {
    wkeOnDocumentReady(window, goOnDocumentReady, param);
}

void goOnWindowDestroy(wkeWebView window, void *param);
void onWindowDestroy(wkeWebView window, void* param) {
    wkeOnWindowDestroy(window, goOnWindowDestroy, param);
}

bool goOnLoadUrlBegin(wkeWebView window, void *param, const utf8* url, wkeNetJob job);
void onLoadUrlBegin(wkeWebView window, void *param) {
    wkeOnLoadUrlBegin(window, goOnLoadUrlBegin, param);
}

void goOnLoadUrlEnd(wkeWebView window, void *param, const utf8* url, wkeNetJob job, void* buf, int len);
void onLoadUrlEnd(wkeWebView window, void *param) {
    wkeOnLoadUrlEnd(window, goOnLoadUrlEnd, param);
}