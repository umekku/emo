# Mind OS API テストコマンド集 (PowerShell)

Write-Host "=== Mind OS API テストコマンド集 ===" -ForegroundColor Cyan
Write-Host ""

# 1. 感覚入力処理
Write-Host "1. 感覚入力処理 (Sensory Input)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sensory" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{"text":"AIFROST、きみのデザイン最高だね！"}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 2. 睡眠処理
Write-Host "2. 睡眠処理 (Sleep)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sleep" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 3. ポジティブフィードバック
Write-Host "3. ポジティブフィードバック (Positive Feedback)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/feedback" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{"positive":true}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 4. ネガティブフィードバック
Write-Host "4. ネガティブフィードバック (Negative Feedback)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/feedback" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{"positive":false}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 5. 脳の状態取得
Write-Host "5. 脳の状態取得 (Get State)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/state" -Method GET -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 6. ストレス適用
Write-Host "6. ストレス適用 (Apply Stress)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/stress" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{"level":80}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 7. 休息適用
Write-Host "7. 休息適用 (Apply Rest)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/rest" -Method POST -ContentType "application/json; charset=utf-8" -Body ''{"quality":100}'' -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

# 8. 直近の記憶取得
Write-Host "8. 直近の記憶取得 (Get Recent Memories)" -ForegroundColor Yellow
Write-Host "コマンド:"
Write-Host 'Invoke-WebRequest -Uri "http://localhost:8081/api/v1/memories" -Method GET -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10' -ForegroundColor Green
Write-Host ""

Write-Host "=== 簡略版（関数として使用） ===" -ForegroundColor Cyan
Write-Host ""

# 関数定義
Write-Host "以下の関数をコピーして使用できます:" -ForegroundColor Yellow
Write-Host ""

$functions = @'
# 感覚入力
function Invoke-MindSensory {
    param([string]$Text)
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sensory" `
        -Method POST `
        -ContentType "application/json; charset=utf-8" `
        -Body "{`"text`":`"$Text`"}" `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# 睡眠
function Invoke-MindSleep {
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/sleep" `
        -Method POST `
        -ContentType "application/json; charset=utf-8" `
        -Body '{}' `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# フィードバック
function Invoke-MindFeedback {
    param([bool]$Positive = $true)
    $body = if ($Positive) { '{"positive":true}' } else { '{"positive":false}' }
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/feedback" `
        -Method POST `
        -ContentType "application/json; charset=utf-8" `
        -Body $body `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# 状態取得
function Get-MindState {
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/state" `
        -Method GET `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# ストレス
function Invoke-MindStress {
    param([int]$Level = 50)
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/stress" `
        -Method POST `
        -ContentType "application/json; charset=utf-8" `
        -Body "{`"level`":$Level}" `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# 休息
function Invoke-MindRest {
    param([int]$Quality = 100)
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/rest" `
        -Method POST `
        -ContentType "application/json; charset=utf-8" `
        -Body "{`"quality`":$Quality}" `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# 記憶取得
function Get-MindMemories {
    Invoke-WebRequest -Uri "http://localhost:8081/api/v1/memories" `
        -Method GET `
        -UseBasicParsing | 
        Select-Object -ExpandProperty Content | 
        ConvertFrom-Json | 
        ConvertTo-Json -Depth 10
}

# 使用例:
# Invoke-MindSensory "AIFROST、きみのデザイン最高だね！"
# Invoke-MindSleep
# Invoke-MindFeedback -Positive $true
# Get-MindState
# Invoke-MindStress -Level 80
# Invoke-MindRest -Quality 100
# Get-MindMemories
'@

Write-Host $functions -ForegroundColor Green
