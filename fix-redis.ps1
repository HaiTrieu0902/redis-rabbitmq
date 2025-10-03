#!/usr/bin/env pwsh
# Fix Redis Connection Script

Write-Host "=== Redis Connection Fix Script ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Check current Redis processes
Write-Host "Step 1: Checking running Redis processes..." -ForegroundColor Yellow
$redisProcesses = netstat -ano | Select-String ":6379" | Select-String "LISTENING"
Write-Host $redisProcesses
Write-Host ""

# Step 2: Find local Redis PIDs (exclude Docker)
Write-Host "Step 2: Finding local Redis processes to stop..." -ForegroundColor Yellow
$localRedisPIDs = @()
$lines = netstat -ano | Select-String ":6379" | Select-String "LISTENING"
foreach ($line in $lines) {
    if ($line -match "127\.0\.0\.1:6379.*LISTENING\s+(\d+)") {
        $pid = $matches[1]
        if ($localRedisPIDs -notcontains $pid) {
            $localRedisPIDs += $pid
            Write-Host "Found local Redis PID: $pid" -ForegroundColor Red
        }
    }
}

# Step 3: Stop local Redis processes
if ($localRedisPIDs.Count -gt 0) {
    Write-Host ""
    Write-Host "Step 3: Stopping local Redis processes..." -ForegroundColor Yellow
    foreach ($pid in $localRedisPIDs) {
        Write-Host "Stopping PID $pid..." -ForegroundColor Red
        try {
            taskkill /F /PID $pid 2>$null
            Write-Host "✅ Stopped PID $pid" -ForegroundColor Green
        } catch {
            Write-Host "❌ Failed to stop PID $pid" -ForegroundColor Red
        }
    }
} else {
    Write-Host "No local Redis processes found to stop." -ForegroundColor Green
}

# Step 4: Check for Redis/Memurai Windows Service
Write-Host ""
Write-Host "Step 4: Checking for Redis Windows Services..." -ForegroundColor Yellow
$redisServices = Get-Service | Where-Object {$_.Name -like "*redis*" -or $_.Name -like "*memurai*"}
if ($redisServices) {
    foreach ($service in $redisServices) {
        Write-Host "Found service: $($service.Name) - Status: $($service.Status)" -ForegroundColor Red
        if ($service.Status -eq "Running") {
            Write-Host "Stopping service: $($service.Name)..." -ForegroundColor Yellow
            try {
                Stop-Service $service.Name -Force
                Write-Host "✅ Stopped service: $($service.Name)" -ForegroundColor Green
            } catch {
                Write-Host "❌ Failed to stop service: $($service.Name)" -ForegroundColor Red
            }
        }
    }
} else {
    Write-Host "No Redis Windows Services found." -ForegroundColor Green
}

# Step 5: Verify only Docker Redis is running
Write-Host ""
Write-Host "Step 5: Verifying only Docker Redis is running..." -ForegroundColor Yellow
Start-Sleep -Seconds 2
$remainingProcesses = netstat -ano | Select-String ":6379" | Select-String "LISTENING"
Write-Host $remainingProcesses
Write-Host ""

# Step 6: Test Docker Redis connection
Write-Host "Step 6: Testing Docker Redis connection..." -ForegroundColor Yellow
try {
    $pingResult = docker exec redis redis-cli PING 2>&1
    if ($pingResult -eq "PONG") {
        Write-Host "✅ Docker Redis is responding!" -ForegroundColor Green
    } else {
        Write-Host "❌ Docker Redis is not responding properly" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Cannot connect to Docker Redis" -ForegroundColor Red
}

# Step 7: Check keys in Docker Redis
Write-Host ""
Write-Host "Step 7: Checking keys in Docker Redis..." -ForegroundColor Yellow
try {
    $keys = docker exec redis redis-cli KEYS "*"
    if ($keys) {
        Write-Host "Keys found:" -ForegroundColor Green
        Write-Host $keys
    } else {
        Write-Host "No keys found (this is OK if you haven't logged in yet)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Cannot query Docker Redis" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== Fix Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "1. Restart your NestJS services (authentication and user)" -ForegroundColor White
Write-Host "2. Login via OAuth (Google/GitHub/MSAL)" -ForegroundColor White
Write-Host "3. Check RedisInsight at http://localhost:8001" -ForegroundColor White
Write-Host "4. You should now see token:* and session:* keys!" -ForegroundColor White
Write-Host ""
