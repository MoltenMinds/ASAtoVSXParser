/*******
* Author: Daniel Azar
* Date: 01/10/2020 
* company: MoltenMinds
********/

package models

/*  Example of an ASA information:
route OUTSIDE 0.0.0.0 0.0.0.0 181.209.69.246 1
route INSIDE_DCINT-MINMOD-DNITO-MIG-511-VPN 10.1.0.0 255.255.224.0 192.168.21.5 1
*/

// Route Contains the information to define a route rule of a switch.
type Route struct {
	Name     string
	NextHop  string
	IPAddres string
	Netmask  string
}
