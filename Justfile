run:
    go build
    ./gh-box

test:
    go test ./...

snap:
    UPDATE_SNAPS=true go test ./...
