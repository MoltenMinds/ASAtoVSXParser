/*******
* Author: Daniel Azar
* Date: 01/10/2020 
* company: MoltenMinds
********/

package models

/*  Example of an ASA interface information:
interface Port-channel1.340
description INSIDE_93_DCINT_MAN_RFG_MINTRANS_EXCH
nameif INSIDE_93_DCINT_MAN_RFG_MINTRANS_EXCH
security-level 30
ip address 10.240.115.98 255.255.255.248
*/

// Port Contains the information to define an interface of a switch.
type Port struct {
	ID            string
	Name          string
	Description   string
	Shutdown      bool
	SecurityLevel int
	IPAddres      string
	Netmask       string
}
