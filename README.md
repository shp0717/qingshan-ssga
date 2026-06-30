# 青山國中小學生自治會官方網站

## 一、概述

我們正在為[青山國中小](https://www.csjhs.ntpc.edu.tw)的[學生自治會](https://ssga.80443.dev)開發一個官方網站。  
該網站使用 [Go](https://golang.org) 語言構建，使用 [JSON](https://www.json.org/json-zh.html) 文件作為數據庫。  
該網站目前正在開發中，我們歡迎社區的貢獻。  
如果您有任何建議或想做出貢獻，歡迎提交[合併請求](https://github.com/shp0717/qingshan-ssga/pulls)或[開啟問題](https://github.com/shp0717/qingshan-ssga/issues)。

## 二、功能與頁面

- 首頁
- 最新消息
- 活動公告
- 意見回饋
- 關於我們
- 聯絡我們
- 開發人員名單

## 三、安裝與運行

### 步驟 1: 下載源代碼

#### 方法1: 使用 Git 克隆倉庫

```bash
git clone https://github.com/shp0717/qingshan-ssga.git
```

#### 方法2: 下載 ZIP 文件

點選[這裡](https://github.com/shp0717/qingshan-ssga/archive/refs/heads/main.zip)下載最新的 ZIP 文件，然後解壓縮到您想要的目錄。

### 步驟 2: 開啟終端機

在您的電腦上開啟終端機（Terminal），並開啟到您下載或解壓縮的目錄。

### 步驟 3: 運行網站

```bash
cd bin
./server-linux-amd64 --host 0.0.0.0 --port 8080  # 如果您使用的是 Linux
server-windows-amd64.exe --host 0.0.0.0 --port 8080  # 如果您使用的是 Windows CMD
.\server-windows-amd64.exe --host 0.0.0.0 --port 8080  # 如果您使用的是 Windows PowerShell
./server-macos-arm64 --host 0.0.0.0 --port 8080  # 如果您使用的是 macOS
```

### 步驟 4: 訪問網站

在您的瀏覽器中輸入以下網址：

```
http://localhost:8080
```

如果您想要在局域網中訪問，請將 `localhost` 替換為您的電腦的局域網 IP 地址。  
如果您想要在互聯網上訪問，請確保您的路由器已經設置了端口轉發，並將 `localhost` 替換為您的公共 IP 地址。  
如果您想要使用你的域名訪問，請確保您的域名已經解析到您的公共 IP 地址。

## 四、貢獻指南

我們歡迎社區的貢獻。如果您有任何建議或想做出貢獻，歡迎提交[合併請求](https://github.com/shp0717/qingshan-ssga/pulls)或[開啟問題](https://github.com/shp0717/qingshan-ssga/issues)。
