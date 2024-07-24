go vet
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags -s