#ifndef WEBVIEW_DEFINE_H
#define WEBVIEW_DEFINE_H

#include "mb.h"
#include "util.h"

void execJs(mbWebView window, const char* code);
mbWebView createWebWindow(int width, int height, bool transparent);
HWND getWindowHandle(mbWebView window);
void setWindowTitle(mbWebView window, const char *title);
void loadURL(mbWebView window, char *url);
void reloadURL(mbWebView window);
void destroyWindow(mbWebView window);
void showWindow(mbWebView window, bool show);
void moveToCenter(mbWebView window);
void setLocalStorageFullPath(mbWebView webView, const char* path);
void setCookieJarFullPath(mbWebView webView, const char* path);
void showDevtools(mbWebView webView, const char* path);

void onDocumentReady(mbWebView window, void* param);
void onWindowDestroy(mbWebView window, void* param);
void onLoadUrlBegin(mbWebView window, void *param);
void onLoadUrlEnd(mbWebView window, void *param);

#endif