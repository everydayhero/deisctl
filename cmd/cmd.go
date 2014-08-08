package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/deis/deisctl/client"
)

func List(c client.Client) error {
	err := c.List()
	return err
}

func Scale(c client.Client, targets []string) error {
	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		err = c.Scale(component, num)
		if err != nil {
			return err
		}
	}
	return nil
}

func Start(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Start(target, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func Stop(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Stop(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Status(c client.Client, targets []string) error {
	for _, target := range targets {
		err := c.Status(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Install(c client.Client) error {
	// data containers
	dataContainers := []string{
		"database-data",
		"registry-data",
		"logger-data",
		"builder-data",
	}
	fmt.Println("Scheduling data containers...")
	for _, dataContainer := range dataContainers {
		c.Create(dataContainer, true)
		// if err != nil {
		// 	return err
		// }
	}
	fmt.Println("Activating data containers...")
	for _, dataContainer := range dataContainers {
		c.Start(dataContainer, true)
		// if err != nil {
		// 	return err
		// }
	}
	// start service containers
	targets := []string{
		"database=1",
		"cache=1",
		"logger=1",
		"registry=1",
		"controller=1",
		"builder=1",
		"router=1"}
	fmt.Println("Scheduling units...")
	err := Scale(c, targets)
	fmt.Println("Activating units...")
	err = Start(c, []string{"registry", "logger", "cache", "database"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"controller"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"builder"})
	if err != nil {
		return err
	}
	err = Start(c, []string{"router"})
	if err != nil {
		return err
	}
	fmt.Println("Done.")
	return err
}

func Uninstall(c client.Client) error {
	targets := []string{
		"database=0",
		"cache=0",
		"logger=0",
		"registry=0",
		"controller=0",
		"builder=0",
		"router=0"}
	err := Scale(c, targets)
	return err
}

func splitScaleTarget(target string) (c string, num int, err error) {
	r := regexp.MustCompile(`([a-z-]+)=([\d]+)`)
	match := r.FindStringSubmatch(target)
	if len(match) == 0 {
		err = fmt.Errorf("Could not parse: %v", target)
		return
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}
