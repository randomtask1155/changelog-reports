export GOPATH=`pwd`
export PATH=$GOPATH/bin:$PATH;

go test gather_changelogs
go test render_changelogs




