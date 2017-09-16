export GOPATH=`pwd`
export PATH=$GOPATH/bin:$PATH;

go get github.com/lib/pq

go build -o build/gather_changelogs gather_changelogs
go build -o build/render_changelogs render_changelogs


GOOS=linux GOARCH=amd64 go build -o build/gather_changelogs_amd64 gather_changelogs
GOOS=linux GOARCH=amd64 go build -o build/render_changelogs_amd64 render_changelogs


