package timod

// Version exposes the version of this library
const Version = "2.0.0"

//Publish module:
//
//  go get -u ./...
//  go mod tidy
//  git tag {VERSION}
//  git push origin {VERSION}
//  GOPROXY=proxy.golang.org go list -m github.com/thingsdb/go-timod@{VERSION}
