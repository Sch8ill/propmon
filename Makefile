bin_name=propmon
target=cmd/main.go

all: build

run:
	go run $(target)

build:
	go build -ldflags="-s" -trimpath -o build/$(bin_name) $(target)

clean:
	rm -rf build