export GOPATH=`pwd`
export PATH=$GOPATH/bin:$PATH;

go build gather_changelogs
go build render_changelogs


GOOS=linux GOARCH=amd64 go build -o gather_changelogs_amd64 gather_changelogs
GOOS=linux GOARCH=amd64 go build -o render_changelogs_amd64 render_changelogs


