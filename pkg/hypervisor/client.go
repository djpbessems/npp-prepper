package hypervisor

import (
	"context"
	neturl "net/url"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

func NewClient(ctx context.Context, host, username, password string, insecure bool) (*vim25.Client, error) {
	url, err := soap.ParseURL(host)
	if err != nil {
		return nil, err
	}

	url.User = neturl.UserPassword(username, password)

	session := &cache.Session{
		URL:      url,
		Insecure: insecure,
	}

	clt := new(vim25.Client)
	err = session.Login(ctx, clt, nil)
	if err != nil {
		return nil, err
	}

	return clt, nil
}

func DatacenterFinder(ctx context.Context, clt *vim25.Client, datacenter string) (*find.Finder, error) {
	fnd := find.NewFinder(clt)

	founddc, err := fnd.DatacenterOrDefault(ctx, datacenter)
	if err != nil {
		return nil, err
	}

	return fnd.SetDatacenter(founddc), nil
}
