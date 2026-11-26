package timod

// Version exposes the version of this library
const Version = "2.0.1"

//Publish module:
//
//  go get -u ./...
//  go mod tidy
//  git tag {VERSION}
//  git push origin {VERSION}
//  GOPROXY=proxy.golang.org go list -m github.com/thingsdb/go-timod/v2@{VERSION}
