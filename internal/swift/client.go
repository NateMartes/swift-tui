package swift

import (
	"context"
	"fmt"
	"strconv"
	"github.com/NateMartes/swift-tui/pkg/errors"
	"github.com/NateMartes/swift-tui/pkg/util"
	"github.com/ncw/swift/v2"
	"strings"
)

// Get arguments for tempauth login as
// hostname, port, username, api-key, use_https
func GetTempAuthArgs() (string, int, string, string, bool) {

	util.LogDebug("Gathering arguments for tempauth")
	username, err := util.UsernameVal()
	if util.UsernameSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if !util.UsernameSupplied() {
		util.LogWarning("No username supplied for tempauth")
	}
	
	password, err := util.ApiKeyVal()
	if util.ApiKeySupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if !util.ApiKeySupplied() {
		util.LogWarning("No api-key/password supplied for tempauth")
	}
	
	swiftHostname, err := util.SwiftHostnameVal()
	if util.SwiftHostnameSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	swiftPort, err := util.SwiftPortVal()
	if util.SwiftPortSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	useHTTPS, err := util.NoHTTPSVal()
	if util.NoHTTPSSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}

	return swiftHostname, swiftPort, username, password, useHTTPS
}

func SetClientFromCloudsFile(filepath string) *swift.Connection {
   
	cloud := GetCloudFromCloudsFile(filepath)

	var method AuthMethod
    switch {
	    case strings.Contains(cloud.Auth.AuthURL, "/v3"):
	        method = AuthKeystoneV3Password
	    case strings.Contains(cloud.Auth.AuthURL, "/v2.0"):
	        method = AuthKeystoneV2
	    default:
	        method = AuthTempauth
    }
    
    client := GetSwiftClientFromClouds(
        context.Background(),
        cloud.Auth.AuthURL,
        cloud.Auth.Username,
        cloud.Auth.Password,
        cloud.Auth.TenantName,
        cloud.Auth.ProjectDomainName,
        method,
    )

    return client
}

func SetClientFromTempauth() *swift.Connection {
	
	hostname, port, username, apiKey, useHttps := GetTempAuthArgs()
	client := GetSwiftClientFromTempauth(context.Background(), hostname, port, username, apiKey, useHttps)
	if client == nil || !client.Authenticated() {
		util.LogFatal("Failed to authenticate client with OpenStack Swift", errors.AUTH_ERROR)
	}
	util.LogDebug(fmt.Sprintf("Authenticated with OpenStack Swift as %s", username))
	return client
	
}

// Gets an OpenStack Swift client using Swift's tempauth middleware
func GetSwiftClientFromTempauth(
	ctx context.Context, 
	hostname string, 
	port int, 
	username string, 
	password string, 
	use_https bool,
) *swift.Connection {
    var url string
    if use_https {
        url = "https://"
    } else {
        url = "http://"
    }
    url += hostname + ":" + strconv.Itoa(port) + "/auth/v1.0"

    util.LogDebug(
    	fmt.Sprintf(
        "Setting auth options for OpenStack Swift tempauth to:\n"+
            "	URL: %s\n"+
            "	Username: %s\n"+
            "	API Key / Password: *****", url, username),
    )

    c := &swift.Connection{
        UserName: username,
        ApiKey:   password,
        AuthUrl:  url,
    }

    if err := c.Authenticate(ctx); err != nil {
        util.LogFatal(fmt.Sprintf("OpenStack Swift tempauth failed: %v", err), errors.AUTH_ERROR)
    }

    return c
}

// Gets an OpenStack Swift client using OpenStack Keystone authentication
func GetSwiftClientFromClouds(
    ctx context.Context,
    url string,
    username string,
    password string,
    tenant string,
    domain string,
    authVersion AuthMethod,
) *swift.Connection {

    util.LogDebug(
        fmt.Sprintf(
            "Setting auth options for OpenStack Swift with Keystone to:\n"+
                "	URL: %s\n"+
                "	Username: %s\n"+
                "	Tenant: %s\n"+
                "	Domain: %s\n"+
                "	Auth Version: %d\n"+
                "	API Key / Password: *****",
            url, username, tenant, domain, authVersion),
    )

    c := &swift.Connection{
        UserName:    username,
        ApiKey:      password,
        AuthUrl:     url,
        Tenant:      tenant,
        Domain:      domain,
        AuthVersion: int(authVersion),
    }

    if err := c.Authenticate(ctx); err != nil {
        util.LogFatal(fmt.Sprintf("OpenStack Swift authentication failed: %v", err), errors.AUTH_ERROR)
    }
    return c
}