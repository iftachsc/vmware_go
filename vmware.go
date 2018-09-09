package vmware

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

// getEnvString returns string from environment variable.
func getEnvString(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	return r
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(v string, def bool) bool {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	switch strings.ToLower(r[0:1]) {
	case "t", "y", "1":
		return true
	}

	return false
}

/*
export GOVMOMI_URL=https://vcenter01.mgmt.il-center-1.cloudzone.io
export GOVMOMI_USERNAME=iftachsc
export GOVMOMI_PASSWORD="5tw5j;M]HVN0$"
export GOVMOMI_INSECURE=true
*/
const (
	envURL      = "GOVMOMI_URL"
	envUserName = "GOVMOMI_USERNAME"
	envPassword = "GOVMOMI_PASSWORD"
	envInsecure = "GOVMOMI_INSECURE"
)

func processOverride(u *url.URL) {
	envUsername := os.Getenv(envUserName)
	envPassword := os.Getenv(envPassword)

	// Override username if provided
	if envUsername != "" {
		var password string
		var ok bool

		if u.User != nil {
			password, ok = u.User.Password()
		}

		if ok {
			u.User = url.UserPassword(envUsername, password)
		} else {
			u.User = url.User(envUsername)
		}
	}

	// Override password if provided
	if envPassword != "" {
		var username string

		if u.User != nil {
			username = u.User.Username()
		}

		u.User = url.UserPassword(username, envPassword)
	}
}

// NewClient creates a govmomi.Client for use in the examples
func NewClientFromEnv(ctx context.Context) (*govmomi.Client, error) {

	var urlDescription = fmt.Sprintf("ESX or vCenter URL [%s]", envURL)
	var urlFlag = flag.String("url", getEnvString(envURL, "https://username:password@host"), urlDescription)

	var insecureDescription = fmt.Sprintf("Don't verify the server's certificate chain [%s]", envInsecure)
	var insecureFlag = flag.Bool("insecure", getEnvBool(envInsecure, false), insecureDescription)

	flag.Parse()

	// Parse URL from string
	u, err := soap.ParseURL(*urlFlag)
	if err != nil {
		return nil, err
	}

	// Override username and/or password as required
	processOverride(u)

	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, *insecureFlag)
}

func NewClient(ctx context.Context, host string, username string, password string) (*govmomi.Client, error) {

	// Parse URL from string
	u, err := soap.ParseURL(host)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(username, envPassword)
	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, true)
}
