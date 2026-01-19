# Mind OS API クイックリファレンス (PowerShell)

## 基本コマンド形式

```powershell
Invoke-WebRequest -Uri "URL" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"key":"value"}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

## API エンドポイント

### 1. 感覚入力処理
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sensory" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"text":"AIFROST、きみのデザイン最高だね！"}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 2. 睡眠処理
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sleep" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 3. ポジティブフィードバック
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/feedback" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"positive":true}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 4. ネガティブフィードバック
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/feedback" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"positive":false}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 5. 脳の状態取得
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/state" `
  -Method GET `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 6. ストレス適用
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/stress" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"level":80}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 7. 休息適用
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/rest" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"quality":100}' `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

### 8. 直近の記憶取得
```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/memories" `
  -Method GET `
  -UseBasicParsing | 
  Select-Object -ExpandProperty Content | 
  ConvertFrom-Json | 
  ConvertTo-Json -Depth 10
```

## 注意事項

### PowerShellでのcurlエイリアス問題

PowerShellでは `curl` は `Invoke-WebRequest` のエイリアスですが、
パラメータ形式が異なるため、以下のような書き方はエラーになります:

```bash
# ❌ これはPowerShellでは動作しません
curl -X POST http://localhost:8081/api/v1/sleep \
  -H "Content-Type: application/json" \
  -d '{}'
```

代わりに、`Invoke-WebRequest` を直接使用してください:

```powershell
# ✅ PowerShellではこちらを使用
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sleep" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{}' `
  -UseBasicParsing
```

### 日本語テキストの扱い

日本語を含むJSONを送信する場合は、必ず `charset=utf-8` を指定してください:

```powershell
-ContentType "application/json; charset=utf-8"
```

### バッククォート (`) の使用

PowerShellでは、長いコマンドを複数行に分割する際に、
バックスラッシュ (`\`) ではなく **バッククォート (`)** を使用します。

```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sensory" `
  -Method POST `
  -ContentType "application/json; charset=utf-8" `
  -Body '{"text":"テスト"}' `
  -UseBasicParsing
```
