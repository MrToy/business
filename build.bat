set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
set target=registry.cn-hangzhou.aliyuncs.com/toy/business
go build -o app ./examples/main.go
if "%errorlevel%" NEQ "0" goto :failed
docker build -t %target% .
docker push %target%
if "%errorlevel%" NEQ "0" goto :failed
goto :end

:failed
echo failed
pause

:end