#include "webview.h"

void execJs(mbWebView webView, const char* code) {
    mbRunJs(webView, mbWebFrameGetMainFrame(webView), code, FALSE, NULL, NULL, NULL);
}

mbWebView createWebWindow(int width, int height, bool transparent) {
    mbWebView webView = mbCreateWebWindow(transparent ? MB_WINDOW_TYPE_TRANSPARENT : MB_WINDOW_TYPE_POPUP, NULL, 0, 0, width, height);
    return webView;
}

HWND getWindowHandle(mbWebView webView)
{
    return mbGetHostHWND(webView);
}

void setWindowTitle(mbWebView webView, const char *title) {
    HWND hWnd = getWindowHandle(webView);
    SetWindowTextW(hWnd, utf8ToUtf16(title));
}

void loadURL(mbWebView webView, char *url)
{
    mbLoadURL(webView, url);
    free(url);
}

void reloadURL(mbWebView webView)
{
    mbReload(webView);
}

void destroyWindow(mbWebView webView)
{
    mbDestroyWebView(webView);
}

void showWindow(mbWebView webView, bool show) {
    mbShowWindow(webView, show);
}

void moveToCenter(mbWebView webView) {
    mbMoveToCenter(webView);
}

void setLocalStorageFullPath(mbWebView webView, const char* path) {
    mbSetLocalStorageFullPath(webView, utf8ToUtf16(path));
}

void setCookieJarFullPath(mbWebView webView, const char* path) {
    mbSetCookieJarFullPath(webView, utf8ToUtf16(path));
}

void showDevtools(mbWebView webView, const char* path) {
    mbSetDebugConfig(webView, "showDevTools", path);
}

// ----------- Callback -------------

void goOnDocumentReady(mbWebView webView, void *param, mbWebFrameHandle frameId);
void onDocumentReady(mbWebView webView, void* param) {
    mbOnDocumentReady(webView, goOnDocumentReady, param);
}

BOOL goOnWindowDestroy(mbWebView webView, void *param, void* unused);
void onWindowDestroy(mbWebView webView, void* param) {
    mbOnDestroy(webView, goOnWindowDestroy, param);
}

BOOL goOnLoadUrlBegin(mbWebView webView, void *param, const char* url, void* job);
void onLoadUrlBegin(mbWebView webView, void *param) {
    mbOnLoadUrlBegin(webView, goOnLoadUrlBegin, param);
}

void goOnLoadUrlEnd(mbWebView webView, void *param, const char* url, void* job, void* buf, int len);
void onLoadUrlEnd(mbWebView webView, void *param) {
    mbOnLoadUrlEnd(webView, goOnLoadUrlEnd, param);
}