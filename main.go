package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/cni-k8s-tool/client"
	"github.com/aws/cni-k8s-tool/nl"
	"github.com/aws/cni-k8s-tool/nl/lib"
	"github.com/urfave/cli"
)

func listFreeIps(c *cli.Context) error {
	ips, err := client.FindFreeIPsAtIndex(0, false)
	if err != nil {
		fmt.Println(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "adapter\tip\t")
	for _, ip := range ips {
		fmt.Fprintf(w, "%v\t%v\t\n",
			ip.Interface.LocalName(),
			ip.IP)
	}
	w.Flush()
	return nil
}

func actionRegistryList(c *cli.Context) error {
	return lib.LockfileRun(func() error {

		reg := &client.Registry{}
		ips, err := reg.List()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ip\t")
		for _, ip := range ips {
			fmt.Fprintf(w, "%v\t\n",
				ip)
		}
		w.Flush()
		return nil
	})
}

func actionIPgc(c *cli.Context) error {
	return lib.LockfileRun(func() error {

		reg := &client.Registry{}
		freeAfter := c.Duration("free-after")
		if freeAfter <= 0*time.Second {
			fmt.Fprintf(os.Stderr,
				"Invalid duration specified. free-after must be > 0 seconds. Got %v. Please specify with --free-after=[time]\n", freeAfter)
			return fmt.Errorf("invalid duration")
		}

		// Insert free-after jitter of 15% of the period
		freeAfter = client.Jitter(freeAfter, 0.15)

		// Invert free-after
		freeAfter *= -1

		ips, err := reg.TrackedBefore(time.Now().Add(freeAfter))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		// grab a list of in-use IPs to sanity check
		assigned, err := nl.GetIPs()
		if err != nil {
			return err
		}

	OUTER:
		for _, ip := range ips {
			// forget IPs that are actually in use and skip over
			for _, assignedIP := range assigned {
				if assignedIP.IPNet.IP.Equal(ip) {
					reg.ForgetIP(ip)
					continue OUTER
				}
			}
			err := client.DefaultClient.DeallocateIP(&ip)
			if err == nil {
				reg.ForgetIP(ip)
			} else {
				fmt.Fprintf(os.Stderr, "Can't deallocate %v due to %v", ip, err)
			}
		}

		return nil
	})
}

func main() {
	var version string
	version = "0.1-beta"

	if os.Getuid() != 0 {
		fmt.Fprintln(os.Stderr, "This command must be run as root")
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "list-free-IPs",
			Usage:  "List all currently unassigned AWS IP addresses",
			Action: listFreeIps,
		},
		{
			Name:   "registry-list",
			Usage:  "List all known free IPs in the internal registry",
			Action: actionRegistryList,
		},
		{
			Name:   "ip-gc",
			Usage:  "Free all IPs that have remained unused for a given time interval",
			Action: actionIPgc,
			Flags: []cli.Flag{
				cli.DurationFlag{Name: "free-after",
					Value: 0 * time.Second},
			},
		},
	}

	app.Version = version
	app.Usage = "Interface with ENI adapters and CNI bindings for those"
	app.Run(os.Args)
}
