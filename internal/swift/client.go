package swift

import (
    "context"
    "fmt"
    "strconv"
    "github.com/NateMartes/go-swift-tui/pkg/errors"
    "github.com/NateMartes/go-swift-tui/pkg/util"
    "github.com/ncw/swift/v2"
)

func GetTempauthClient(ctx context.Context, hostname string, port int, username string, password string, use_https bool) *swift.Connection {
    var url string
    if use_https {
        url = "https://"
    } else {
        url = "http://"
    }
    url += hostname + ":" + strconv.Itoa(port) + "/auth/v1.0"

    util.LogDebug(fmt.Sprintf(
        "Setting auth options for OpenStack Swift tempauth to:\n"+
            "	URL: %s\n"+
            "	Username: %s\n"+
            "	Password: *****", url, username))

    c := &swift.Connection{
        UserName: username,
        ApiKey:   password,
        AuthUrl:  url,
    }

    if err := c.Authenticate(ctx); err != nil {
        util.LogFatal(fmt.Sprintf("OpenStack Swift tempauth failed: %v\n", err), errors.AUTH_ERROR)
    }

    return c
}