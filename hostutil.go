package vmware

import (
	"context"
	"log"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func GetEsxHost(c *govmomi.Client, ctx context.Context) ([]mo.HostSystem, error) {
	//defer m.Destroy(ctx)

	//defer c.Logout(ctx)
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)

	var hostsystems []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hostsystems)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return hostsystems, nil
}
