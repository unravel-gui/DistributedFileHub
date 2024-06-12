@echo off
setlocal enabledelayedexpansion

REM 设置循环次数
set /a count=2

REM 执行循环
for /l %%i in (1,1,%count%) do (
    echo Running iteration %%i
    start cmd /k "go run ..\usernode\usernode_server.go -f=..\config\usernode%%i.json"
)

echo All iterations completed.