/*******
* Author: Daniel Azar
* Date: 01/10/2020 
* company: MoltenMinds
********/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MoltenMinds/ASAtoVSXParser/models"
)

var filename = "input"
var baseName = "input.txt"
var outputFile = "-output"
var interfacesJSON = "-detail-port"
var routeJSON = "-detail-route"
var extension = ".txt"

/*
* Main method, open files, then scans them for the configuration information and generates output files.
*/
func run() {

	if len(os.Args) > 1 {
		baseName = os.Args[1]
		filename = strings.TrimPrefix(baseName, ".\\")
		extension = filepath.Ext(filename)
		filename = strings.TrimRight(filename, extension)
		if len(os.Args) > 2 {
			outputFile = os.Args[2]
		} else {
			outputFile = filename + outputFile + extension
		}
		routeJSON = filename + routeJSON + extension
		interfacesJSON = filename + interfacesJSON + extension
	} else {
		outputFile = filename + outputFile + extension
	}
	file, err := os.Open(baseName)
	if err != nil {
		log.Fatalf("failed opening input file: %s", err)
	}
	defer file.Close()
	routeJSON = filename + routeJSON + extension
	interfacesJSON = filename + interfacesJSON + extension
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	var ports []models.Port
	var routes []models.Route

	/*
	Example of expected ASA file:
	interface Port-channel1.340
		  description INSIDE_93_DCINT_MAN_RFG_MINTRANS_EXCH
		  nameif INSIDE_93_DCINT_MAN_RFG_MINTRANS_EXCH
		  security-level 30
		  ip address 10.240.115.98 255.255.255.248

	route OUTSIDE 0.0.0.0 0.0.0.0 181.209.69.246 1
	route INSIDE_DCINT-MINMOD-DNITO-MIG-511-VPN 10.1.0.0 255.255.224.0 192.168.21.5 1
	*/

	for scanner.Scan() {
		var port models.Port
		var route models.Route
		line := scanner.Text()
		txtlines = append(txtlines, line)

		var _, err = fmt.Sscanf(line, "interface Port-channel%s", &port.ID)

		if err == nil {
			scanner.Scan()
			line = scanner.Text()
			_, err := fmt.Sscanf(line, " description %s", &port.Description)
			if err == nil {
				scanner.Scan()
				line = scanner.Text()
			}
			if line == " shutdown" {
				port.Shutdown = true
				scanner.Scan()
				line = scanner.Text()
			}
			fmt.Sscanf(line, " nameif %s", &port.Name)
			scanner.Scan()
			line = scanner.Text()
			fmt.Sscanf(line, " security-level %d", &port.SecurityLevel)
			scanner.Scan()
			line = scanner.Text()
			fmt.Sscanf(line, " ip address %s %s", &port.IPAddres, &port.Netmask)
			ports = append(ports, port)
		}
		_, err = fmt.Sscanf(line, "route %s %s %s %s", &route.Name, &route.IPAddres, &route.Netmask, &route.NextHop)
		if err == nil {
			routes = append(routes, route)
		}
	}
	fmt.Printf("Number of Interface: %d\n", len(ports))
	fmt.Printf("Number of Routes: %d\n", len(routes))
	writeFile(ports, routes)
	writeInterfacesJSON(ports)
	writeRoutesJSON(routes)

}

/**
* Function: writeFile
* Writes the output file containing the VSX configuration commands.
* The bond%s in the Fprintf could be renamed for eth%s or other convenient name.
*/
func writeFile(ports []models.Port, routes []models.Route) {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed opening output file: %s", err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "transaction begin\n")

	for _, port := range ports {
		if port.Shutdown {
			fmt.Fprintf(writer, "#add interface name bond%s ip %s netmask %s\n", port.ID, port.IPAddres, port.Netmask)
		} else {
			fmt.Fprintf(writer, "add interface name bond%s ip %s netmask %s\n", port.ID, port.IPAddres, port.Netmask)
		}
	}
	for _, route := range routes {
		if route.IPAddres == "0.0.0.0" {
			fmt.Fprintf(writer, "add route destination default next_hop %s\n", route.NextHop)
		} else {
			fmt.Fprintf(writer, "add route destination %s netmask %s next_hop %s\n", route.IPAddres, route.Netmask, route.NextHop)
		}
	}

	fmt.Fprintf(writer, "transaction end\n")
	writer.Flush()
}

/**
* Function: writeInterfacesJSON
* Writes the interfaces model array received to a JSON file.
*/
func writeInterfacesJSON(ports []models.Port) {
	file, err := os.Create(interfacesJSON)
	if err != nil {
		log.Fatalf("failed opening output file: %s", err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	json, err := json.MarshalIndent(ports, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(writer, string(json))
	writer.Flush()
}

/**
* Function: writeRoutesJSON
* Writes the routes model array received to a JSON file.
*/
func writeRoutesJSON(routes []models.Route) {
	file, err := os.Create(routeJSON)
	if err != nil {
		log.Fatalf("failed opening output file: %s", err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	json, err := json.MarshalIndent(routes, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(writer, string(json))
	writer.Flush()
}

func main() {
	fmt.Println(("Program started"))
	run()
	fmt.Println(("Program completed"))

}
