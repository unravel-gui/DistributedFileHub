@echo off
setlocal enabledelayedexpansion

REM 设置循环次数
set /a count=7

REM 执行循环
for /l %%i in (1,1,%count%) do (
    echo Running iteration %%i
    go run ..\datanode\datanode_server.go -f=..\config\data_server%%i.json -c
)

echo All iterations completed.