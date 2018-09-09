package vmware

import (
	"context"
	"log"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func GetVM(c *govmomi.Client, ctx context.Context) ([]mo.VirtualMachine, error) {

	//defer c.Logout(ctx)

	// Create view of VirtualMachine objects
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var virtualmachines []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"guest", "config"}, &virtualmachines)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return virtualmachines, nil
}

// foreach ($VirtualSCSIController in ($vmView.Config.Hardware.Device | where {$_.DeviceInfo.Label -match "SCSI Controller"}))
// 	{
// 		foreach ($hd in ($VMView.Config.Hardware.Device | where {$_.ControllerKey -eq $VirtualSCSIController.Key}))
// 		{
// 			$diskNode = $VmDisksNode.AppendChild((createElement -stringId "VMDK"))

// 			if(!$hd.Backing.Datastore)
// 			{
// 				# skipping if hd is not accessible
// 				continue
// 			}

// 			$datastoreId = $hd.Backing.Datastore.GetHashCode()
// 			# csv null
//             $lunMap = $null
// 			$lunPath = $null
//             $lunTarget = $null
// 			$lunVserver = $null
// 			$diskNode.SetAttribute("VmName",$vmView.Name)
// 			if($hd.Backing.GetType().Name.Contains("RawDiskMapping"))
// 			{
// 				$diskNode.SetAttribute("Type","RawDiskMapping")

// 				# imatating sdk scsi lun id value
// 				$matchScsiLun = $script:scsiLuns | ?{$_.ExtensionData.Uuid -eq $hd.Backing.LunUuid} | Select -Unique
//                 $scsiLunPath = (Get-ScsiLunPath -ScsiLun $matchScsiLun | ?{$_.State -eq "Active"} | Select -First 1)
//                 $scsiTargetSanId = $scsiLunPath.SanID

//                 if($matchScsiLun.RuntimeName -eq $null)
//                 {
//                     $lunMapId = $scsiLunPath.LunPath.Split(":")[3].Trim('L')
// 				    $hbaDevice = $scsiLunPath.LunPath.Split(":")[0]
//                 }
//                 else
//                 {
//                     $lunMapId = $matchScsiLun.RuntimeName.Split(":")[3].Trim('L')
//                     $hbaDevice = $matchScsiLun.RuntimeName.Split(":")[0]
//                 }
//                 $hba = Get-VMHostHba -VMHost $matchScsiLun.VMHost -Device $hbaDevice

// 				if($hba.IScsiName)
// 				{
//                     foreach ($connection in $global:NC_Connections)
// 					{
// 					    $lunMap  = getIscsiLunPathByRdmDisk -hdExtensionData $hd -connection $connection -hba $hba -lunId $lunMapId
//                         if($lunMap -ne $null -and (Get-NcIscsiNodeName -Controller $connection -VserverContext $lunMap.Vserver | %{$_.Name}) -contains $scsiTargetSanId)

//                         {
// 							$lunPath = $lunMap.Path
//                             $lunTarget = $connection.Name
// 							$lunVserver = $lunMap.Vserver
// 							write-host ("Found NetApp info for RDM scsi lun [Hba: {0}, ID: {1}, Target: {2} on controller {3}" -f $hbaDevice,$lunMapId,$scsiTargetSanId,$connection.name) -foregroundcolor green
// 							break;
//                         }
// 						else
// 						{
// 							write-warning ("Could not find NetApp info for RDM scsi lun [Hba: {0}, ID: {1}, Target: {2} on controller {3}" -f $hbaDevice,$lunMapId,$scsiTargetSanId,$connection.name)
// 						}

//                     }
//                     if($lunMap -eq $null) # if not found on nc contollers - check on na
//                     {
//                         foreach ($connection in $global:NA_Connections)
// 					    {
// 	                        $lunMapId = $matchScsiLun.RuntimeName.Split(":")[3].Trim('L')
//                             $lunMap = Get-NaLunMapByInitiator -Initiator $hba.IScsiName -Controller $connection | ?{$_.LunId -eq $lunMapId}
//                             if($lunMap -ne $null -and (Get-NAIscsiNodeName -Controller $connection) -contains $scsiTargetSanId)
//                             {
// 							    $lunPath = $lunMap.Path
//                                 $lunTarget = $connection.Name
// 								write-host ("Found NetApp info for RDM scsi lun [Hba: {0}, ID: {1}, Target: {2} on controller {3}" -f $hbaDevice,$lunMapId,$scsiTargetSanId,$connection.name) -foregroundcolor green
// 							    break;
//                             }
// 							else
// 							{
// 								write-warning ("Could not find NetApp info for RDM scsi lun [Hba: {0}, ID: {1}, Target: {2} on controller {3}" -f $hbaDevice,$lunMapId,$scsiTargetSanId,$connection.name)
// 							}
//                         }
//                     }
// 				}
// 				else
// 				{
// 					$fiberPortWWN = [Convert]::ToString($hba.PortWorldWideName, 16) -replace '(..(?!$))','$1:'

// 					foreach ($connection in $global:NA_Connections)
// 					{
// 					 	$lunMap = Get-NaLunMapByInitiator -Initiator $fiberPortWWN -Controller $connection | ?{$_.LunId -eq $lunMapId}
// 					 	if($lunMap -ne $null -and $lunPath -eq $null -and (Get-NaFcpInterface -Controller $connection | %{$_.PortName.ToUpper()}) -contains $scsiTargetSanId.toUpper())
// 						{
//                             $lunTarget = $connection.Name
// 							$lunPath = $lunMap.Path
// 						}
// 						else
// 						{
// 							write-warning ("Could not find NetApp info for RDM fibre lun [Hba: {0}, ID: {1}, Target: {2} on controller {3}" -f $hbaDevice,$lunMapId,$scsiTargetSanId,$connection.name)
// 						}
// 					}

// 					foreach ($connection in $global:NC_Connections)
// 					{
// 					 	$lunMap = Get-NcLunMapByInitiator -Initiator $fiberPortWWN -Controller $connection | ?{$_.LunId -eq $lunMapId}
// 					 	if($lunMap -ne $null -and $lunPath -eq $null) {
//                             if((Get-NcFcpInterface -Controller $connection | %{$_.PortName.ToUpper()}) -contains $scsiTargetSanId.toUpper())
//                             {
//                                 $lunTarget = $ncConnectionToClusterName[$connection.Name]
// 							    $lunPath = $lunMap.Path
// 								$lunVserver = $lunMap.Vserver
//                             }
// 						}
// 					}
// 				}
// 			}
// 			else
// 			{
// 				$diskNode.SetAttribute("Type","Flat")
// 			}

// 			[string]$diskFileName = $hd.Backing.FileName
// 			$capacityMB = $hd.CapacityInKB / 1KB

// 			# getting all files associated with the curred disk with file from the LayoutEx.Disk array
// 			$vmLayoutExDisk = $vmView.LayoutEx.Disk | ?{$_.Key -eq $hd.Key}
// 			# recording only files associated with the disk
// 			[array]$arrLayoutExDiskFileKeys = $vmLayoutExDisk.Chain | ?{$_ -is [VMware.Vim.VirtualMachineFileLayoutExDiskUnit]}

// 			# calculating actual size of disk by measuring  all disk file actual size and calculate sum
// 			$sizeOnDatastoreBytes = ($arrLayoutExDiskFileKeys | %{
// 				$_.FileKey} | %{
// 					$intFileKey = $_
// 					# matching the file from the LayoutEx.File tree with matching key file and that represent
// 					# a file that is a diskExtent - part of the disk
// 					$vmView.LayoutEx.File | ?{($_.Key -eq $intFileKey) -and ($_.Type -eq "diskExtent")}
// 				} | Measure-Object -Sum Size).Sum

// 			$sizeOnDatastoreMB = [Math]::Round($sizeOnDatastoreBytes / 1MB, 1)

//             $diskNode.AppendChild((createElement -stringId "SCSI_Location" -InnerText ($VirtualSCSIController.BusNumber.ToString() + ":" + $hd.UnitNumber.ToString())))
//             $diskNode.AppendChild((createElement -stringId "Label" -InnerText $hd.DeviceInfo.Label))
//             $diskNode.AppendChild((createElement -stringId "FilePath" -InnerText $diskFileName))
// 			if($lunPath -and $lunTarget)
//             {
//                 $diskNode.AppendChild((createElement -stringId "LunPath" -InnerText $lunPath))
//                 $diskNode.AppendChild((createElement -stringId "LunTarget" -InnerText $lunTarget))
//                 $diskNode.AppendChild((createElement -stringId "LunVserver" -InnerText $lunVserver))
//             }

// 			$diskNode.AppendChild((createElement -stringId "CapacityMB" -InnerText $capacityMB))
// 			$diskNode.AppendChild((createElement -stringId "SinzeOnDatastoreMB" -InnerText $sizeOnDatastoreMB))
// 			$diskNode.AppendChild((createElement -stringId "datatstore" -InnerText $hd.Backing.Datastore.Name))
//         }

/*
//should go to scenario
func RegisterVM() {
	vm := new(mo.HostStorageSystem)

	canBePoweredOn := sUnmirroredVmdkVmNames.Contains(selectedMachineData.Name)
	//AddToLog(string.Format("VM {0} has a VMDK that is not mirrored.", selectedMachineData.Name), "WARNING");

	//newMachineName = PrefixTB.Text + selectedMachineData.Name;
	//newVmxPath = MyVMWare.GetVmxPath(selectedMachineData.VMXPath, PrefixTB.Text);

	//HostSystem esx = null;
	//var replicatedRdms = WizardData.RdmVolumeSnapshotMappingCmode.Keys.Where(x => x.rdmVmdk.vmName == newMachineName.Substring(prefix.Length));

}
*/
