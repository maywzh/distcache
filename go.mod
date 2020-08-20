module main

go 1.14

replace geecache => /Users/maywzh/Workspace/DistCache/geecache //本地包相对路径或绝对路径

require (
	geecache v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.2 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
