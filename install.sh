go build -o aswsg

GOOS=linux GOARCH=amd64   go build -o precompiled/aswsg_LinuxIntel64
GOOS=linux GOARCH=386     go build -o precompiled/aswsg_linuxIntel32
GOOS=linux GOARCH=arm64   go build -o precompiled/aswsg_linuxArm64
GOOS=linux GOARCH=arm     go build -o precompiled/aswsg_linuxArm32
GOOS=darwin GOARCH=amd64  go build -o precompiled/aswsg_MacIntel64
GOOS=windows GOARCH=amd64 go build -o precompiled/aswsg_win64.exe
GOOS=windows GOARCH=386   go build -o precompiled/aswsg_win32.exe
