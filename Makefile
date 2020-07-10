all: test build
build:
	go build -ldflags "-s -w" neckless.adviser.com/cmd/neckless

test:
	go test neckless.adviser.com/key
	go test neckless.adviser.com/symmetric
	go test neckless.adviser.com/asymmetric
	go test neckless.adviser.com/casket
	go test neckless.adviser.com/member
	go test neckless.adviser.com/pearl
	go test neckless.adviser.com/kvpearl
	go test neckless.adviser.com/gem
	go test neckless.adviser.com/necklace
	go test neckless.adviser.com/cmd/neckless


