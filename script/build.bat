::rd /s/q release
::md release
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
cd ..
go build
::COPY go-chatroom release\