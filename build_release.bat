cd tweet
go vet
staticcheck .
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
cd ..\cmd\tweet-bsky
go vet
staticcheck .
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\tweet-changelog
go vet
staticcheck .
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\tweetee
go vet
staticcheck .
for %%j in (shadow, defers) do for /f %%i in ('where %%j') do go vet --vettool=%%i
go build -ldflags "-s -H=windowsgui"
cd ..\..