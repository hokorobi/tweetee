cd tweet
go vet
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
cd ..\cmd\tweet-bsky
go vet
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\tweet-changelog
go vet
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\tweetee
go vet
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\..