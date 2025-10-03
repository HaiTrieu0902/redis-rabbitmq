# Start all microservices

Write-Host "========================================"  -ForegroundColor Green
Write-Host "  Starting Redis MQ Microservices"  -ForegroundColor Green
Write-Host "========================================"  -ForegroundColor Green

# Check if Redis is running
Write-Host "`nChecking Redis..." -ForegroundColor Yellow
$redisRunning = docker ps | Select-String "redis"
if ($redisRunning) {
    Write-Host "✓ Redis is running" -ForegroundColor Green
} else {
    Write-Host "✗ Redis is not running" -ForegroundColor Red
    Write-Host "Starting Redis..." -ForegroundColor Yellow
    docker-compose up -d
    Start-Sleep -Seconds 3
    Write-Host "✓ Redis started" -ForegroundColor Green
}

# Create logs directory if not exists
if (!(Test-Path -Path "logs")) {
    New-Item -ItemType Directory -Path "logs" | Out-Null
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  All systems ready!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host "`nStarting services..." -ForegroundColor Yellow
Write-Host ""
Write-Host "Authentication Service: http://localhost:3001" -ForegroundColor Green
Write-Host "User Service: http://localhost:3002" -ForegroundColor Green
Write-Host ""
Write-Host "Press Ctrl+C to stop all services" -ForegroundColor Yellow
Write-Host ""

# Start Authentication Service
Write-Host "Starting Authentication Service..." -ForegroundColor Cyan
$authJob = Start-Job -ScriptBlock {
    Set-Location $args[0]
    npm run start:dev
} -ArgumentList "$PSScriptRoot\authenication"

# Start User Service
Write-Host "Starting User Service..." -ForegroundColor Cyan
$userJob = Start-Job -ScriptBlock {
    Set-Location $args[0]
    npm run start:dev
} -ArgumentList "$PSScriptRoot\user"

# Wait for services to initialize
Write-Host "`nWaiting for services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host "`n✓ Services started!" -ForegroundColor Green
Write-Host ""
Write-Host "View live logs in separate terminals:" -ForegroundColor Cyan
Write-Host "  cd authenication && npm run start:dev" -ForegroundColor White
Write-Host "  cd user && npm run start:dev" -ForegroundColor White
Write-Host ""
Write-Host "Or check job outputs:" -ForegroundColor Cyan
Write-Host "  Receive-Job -Id $($authJob.Id) -Keep" -ForegroundColor White
Write-Host "  Receive-Job -Id $($userJob.Id) -Keep" -ForegroundColor White
Write-Host ""

# Keep script running and show job status
try {
    while ($true) {
        Start-Sleep -Seconds 5
        
        $authState = (Get-Job -Id $authJob.Id).State
        $userState = (Get-Job -Id $userJob.Id).State
        
        if ($authState -ne "Running") {
            Write-Host "`n✗ Authentication Service stopped" -ForegroundColor Red
            break
        }
        if ($userState -ne "Running") {
            Write-Host "`n✗ User Service stopped" -ForegroundColor Red
            break
        }
    }
}
finally {
    Write-Host "`nStopping services..." -ForegroundColor Yellow
    Stop-Job -Id $authJob.Id, $userJob.Id
    Remove-Job -Id $authJob.Id, $userJob.Id
    Write-Host "✓ Services stopped" -ForegroundColor Green
}
