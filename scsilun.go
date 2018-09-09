package vmware

import (
	"context"
	"log"
	"time"

	"github.com/vmware/govmomi"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func GetScsiLunDisks(c *govmomi.Client, ctx context.Context) (map[string]types.ScsiLun, error) {

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)
	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var hostSystems []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, nil, &hostSystems)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	luns := []types.ScsiLun{}

	for _, hostSystem := range hostSystems {

		host := object.NewHostSystem(c.Client, hostSystem.Reference())

		ss, err := host.ConfigManager().StorageSystem(ctx)
		if err != nil {
			return nil, err
		}

		var hss mo.HostStorageSystem
		err = ss.Properties(ctx, ss.Reference(), nil, &hss)
		if err != nil {
			return nil, err
		}
		println(time.Now().String())
		luns = append(luns, getScsiLunDisks(hss)...)
	}

	return uniqueLuns(luns), err
}

func uniqueLuns(luns []types.ScsiLun) map[string]types.ScsiLun {
	//u := []types.ScsiLun{}
	m := make(map[string]types.ScsiLun)

	for _, lun := range luns {
		if _, ok := m[lun.Uuid]; !ok {
			m[lun.Uuid] = lun
			//u = append(u, lun)
		}
	}

	//return u
	return m
}

func getScsiLunDisks(hss mo.HostStorageSystem) (diskLuns []types.ScsiLun) {
	for _, lun := range hss.StorageDeviceInfo.ScsiLun {
		if lun.GetScsiLun().LunType == "disk" {
			//fmt.Println("uuid", lun.GetScsiLun().Uuid, "device name", lun.GetScsiLun().DeviceName, "model", lun.GetScsiLun().Model,
			//	"cacnonical", lun.GetScsiLun().CanonicalName, "vendor", lun.GetScsiLun().Vendor)
			diskLuns = append(diskLuns, *lun.GetScsiLun())
		}
	}
	return diskLuns
}
