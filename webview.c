#include "webview.h"

void execJs(mbWebView webView, const char* code) {
    mbRunJs(webView, mbWebFrameGetMainFrame(webView), code, FALSE, NULL, NULL, NULL);
}

mbWebView createWebWindow(int width, int height, bool transparent) {
    mbWebView window = mbCreateWebWindow(transparent ? MB_WINDOW_TYPE_TRANSPARENT : MB_WINDOW_TYPE_POPUP, NULL, 0, 0, width, height);
    return window;
}

HWND getWindowHandle(mbWebView window)
{
    return mbGetHostHWND(window);
}

void setWindowTitle(mbWebView window, const char *title) {
    HWND hWnd = getWindowHandle(window);
    SetWindowTextW(hWnd, utf8ToUtf16(title));
}

void loadURL(mbWebView window, char *url)
{
    mbLoadURL(window, url);
    free(url);
}

void reloadURL(mbWebView window)
{
    mbReload(window);
}

void destroyWindow(mbWebView window)
{
    mbDestroyWebView(window);
}

void showWindow(mbWebView window, bool show) {
    mbShowWindow(window, show);
}

void moveToCenter(mbWebView window) {
    mbMoveToCenter(window);
}

void setLocalStorageFullPath(mbWebView webView, const char* path) {
    mbSetLocalStorageFullPath(webView, utf8ToUtf16(path));
}

void setCookieJarFullPath(mbWebView webView, const char* path) {
    mbSetCookieJarFullPath(webView, utf8ToUtf16(path));
}

void showDevtools(mbWebView webView, const char* path) {
    // mbShowDevtools(webView, utf8ToUtf16(path), NULL, NULL);
}

// ----------- Callback -------------

void goOnDocumentReady(mbWebView window, void *param, mbWebFrameHandle frameId);
void onDocumentReady(mbWebView window, void* param) {
    mbOnDocumentReady(window, goOnDocumentReady, param);
}

BOOL goOnWindowDestroy(mbWebView window, void *param, void* unused);
void onWindowDestroy(mbWebView window, void* param) {
    mbOnDestroy(window, goOnWindowDestroy, param);
}

BOOL goOnLoadUrlBegin(mbWebView window, void *param, const char* url, void* job);
void onLoadUrlBegin(mbWebView window, void *param) {
    mbOnLoadUrlBegin(window, goOnLoadUrlBegin, param);
}

void goOnLoadUrlEnd(mbWebView window, void *param, const char* url, void* job, void* buf, int len);
void onLoadUrlEnd(mbWebView window, void *param) {
    mbOnLoadUrlEnd(window, goOnLoadUrlEnd, param);
}