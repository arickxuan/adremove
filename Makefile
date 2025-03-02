NAME = adremove

lib_install:
	go install -v github.com/sagernet/gomobile/cmd/gomobile@v0.1.4
	go install -v github.com/sagernet/gomobile/cmd/gobind@v0.1.4


build_android:
	GOOS=android GOARCH=arm64 go build -o adremove .
build_android_lib:
# gomobile bind -v -target android -androidapi 21 -javapkg=top.arick -libname=adlib -trimpath -buildvcs=false -ldflags -s -w -buildid=  ./adcore
	gomobile bind -v -target android -androidapi 21 -javapkg=top.arick -libname=adlib -trimpath -buildvcs=false -o build/adlib.aar ./adcore