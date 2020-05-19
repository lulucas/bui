#ifndef WEBVIEW_DEFINE_H
#define WEBVIEW_DEFINE_H

#include "mb.h"
#include "util.h"

void execJs(mbWebView webView, const char* code);
mbWebView createWebWindow(int width, int height, bool transparent);
HWND getWindowHandle(mbWebView webView);
void setWindowTitle(mbWebView webView, const char *title);
void loadURL(mbWebView webView, char *url);
void reloadURL(mbWebView webView);
void destroyWindow(mbWebView webView);
void showWindow(mbWebView webView, bool show);
void moveToCenter(mbWebView webView);
void setLocalStorageFullPath(mbWebView webView, const char* path);
void setCookieJarFullPath(mbWebView webView, const char* path);
void showDevtools(mbWebView webView, const char* path);

void onDocumentReady(mbWebView webView, void* param);
void onWindowDestroy(mbWebView webView, void* param);
void onLoadUrlBegin(mbWebView webView, void *param);
void onLoadUrlEnd(mbWebView webView, void *param);

#endif