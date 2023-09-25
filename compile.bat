for /F %%i in ('git describe --long') do ( set commitid=%%i)
set flags="-X eGame-demo-back-office-api/cmd/cli/version.version=%commitid%"
go build -ldflags %flags% .\cmd\ginadmin