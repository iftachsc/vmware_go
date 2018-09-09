package vmware

/*
func GetDatastores(c *govmomi.Client, ctx context.Context) ([]types.ScsiLun, error) {
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)
	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var datastores []mo.Datastore
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{"Summary"}, &datastores)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, datastore := range datastores {
		datastore.Info.GetDatastoreInfo().
	}

	return uniqueLuns(luns), err
}
*/
