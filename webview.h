#ifndef WEBVIEW_DEFINE_H
#define WEBVIEW_DEFINE_H

#include "wke.h"

void bindPort();
wkeWebView createWebWindow(int width, int height, bool transparent);
HWND getWindowHandle(wkeWebView window);
void loadURL(wkeWebView window, char *url);
void reloadURL(wkeWebView window);
void setWindowTitle(wkeWebView window, char *title);
void destroyWindow(wkeWebView window);
void showWindow(wkeWebView window, bool show);
void moveToCenter(wkeWebView window);
void setLocalStorageFullPath(wkeWebView webView, const char* path);
void setCookieJarFullPath(wkeWebView webView, const char* path);
void showDevtools(wkeWebView webView, const char* path);

void onDocumentReady(wkeWebView window, void* param);
void onWindowDestroy(wkeWebView window, void* param);
void onLoadUrlBegin(wkeWebView window, void *param);
void onLoadUrlEnd(wkeWebView window, void *param);

#endif