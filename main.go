package main

/*
	/etc/udev/rules.d/i2c-permissions.rules: ACTION=="add", KERNEL=="i2c-[0-1]*", GROUP="ubuntu"
*/

import (
	"fmt"
	"time"

	"github.com/murdinc/mvrdpi/display"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	disp := display.NewDisplay(1, 0x3c, 64, 128)
	disp.MVRDLogo()

	startTime := time.Now().Round(time.Second)
	hostInfo, _ := host.Info()

	for {
		cpuPercent, _ := cpu.Percent(0, false)
		swapMemory, _ := mem.SwapMemory()
		virtualMemory, _ := mem.VirtualMemory()
		uptime, _ := host.Uptime()

		disp.WriteText("MVRDPi v1.0", 1)
		disp.WriteText(fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion), 129)
		// 257
		disp.WriteText(fmt.Sprintf("   Uptime: %s        ", startTime.Sub(time.Unix(startTime.Unix()-int64(uptime), 0)).String()), 385)
		disp.WriteText(fmt.Sprintf("      CPU: %.2F %%   ", cpuPercent[0]), 513)
		disp.WriteText(fmt.Sprintf("  Swp Mem: %.2F %%   ", swapMemory.UsedPercent), 641)
		disp.WriteText(fmt.Sprintf("  Vrt Mem: %.2F %%   ", virtualMemory.UsedPercent), 769)
		// 897
		time.Sleep(time.Second)
	}
}
