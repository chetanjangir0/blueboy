package main

//todo
// wifi on off function
// use ticks to update frequently
// s for scan
// make it show security and network strength of devices
// styles like loading spinner

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/ui"
	// "os"
)

func main() {
	// debug file
	// f, err := os.Create("debug.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	program := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}





// func fetchPaired() tea.Cmd {
// 	return func() tea.Msg {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
//
// 		output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "NAME,TYPE,UUID", "connection", "show").CombinedOutput()
// 		if ctx.Err() == context.DeadlineExceeded {
// 			return fetchPairedMsg{err: fmt.Errorf("timeout")}
// 		}
// 		if err != nil {
// 			return fetchPairedMsg{err: err}
// 		}
// 		outputString := strings.Trim(string(output), "\n")
// 		outputStringSlice := strings.Split(outputString, "\n")
//
// 		var devices = make([]Device, len(outputStringSlice))
// 		for i, d := range outputStringSlice {
// 			deviceInfo := strings.Split(d, ":") // "Device <name> <type>"
// 			if len(deviceInfo) != 3 {
// 				return fetchPairedMsg{err: fmt.Errorf("unexpected number of fields")}
// 			}
// 			devices[i] = Device{
// 				Name: deviceInfo[0],
// 				Type: deviceInfo[1],
// 				UUID: deviceInfo[2],
// 			}
// 		}
// 		return fetchPairedMsg{devices: devices}
// 	}
// }
//
// func startScan() tea.Cmd {
// 	return func() tea.Msg {
// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()
// 		output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "SSID,SECURITY", "device", "wifi", "list").CombinedOutput()
// 		if ctx.Err() == context.DeadlineExceeded {
// 			return startScanMsg{err: fmt.Errorf("Error: connection timed out")}
//
// 		}
// 		if err != nil {
// 			return startScanMsg{err: err}
// 		}
// 		outputString := strings.Trim(string(output), "\n")
// 		outputStringSlice := strings.Split(outputString, "\n")
//
// 		var devices = make([]Device, len(outputStringSlice))
// 		for i, d := range outputStringSlice {
// 			deviceInfo := strings.Split(d, ":")
//
// 			if len(deviceInfo) == 0 { // if there is no ssid/name
// 				continue
// 			} else if len(deviceInfo) == 1 { // if there is no security
// 				deviceInfo = append(deviceInfo, " ")
// 			}
//
// 			devices[i] = Device{
// 				Name:     deviceInfo[0],
// 				Type:     "wifi",
// 				Security: deviceInfo[1],
// 			}
// 		}
// 		return startScanMsg{devices: devices}
// 	}
//
// }
//
// func connectDevice(UUID string) tea.Cmd {
// 	return func() tea.Msg {
// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()
//
// 		_, err := exec.CommandContext(ctx, "nmcli", "connection", "up", UUID).CombinedOutput()
// 		if ctx.Err() == context.DeadlineExceeded {
// 			return connectDeviceMsg{err: fmt.Errorf("Error: connection timed out")}
// 		}
// 		if err != nil {
// 			return connectDeviceMsg{err: err}
// 		}
// 		return connectDeviceMsg{output: "connection successfully activated"}
// 	}
// }
//
// func pairNewDevice(newDevice Device, password string) tea.Cmd {
// 	return func() tea.Msg {
// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()
//
// 		var err error
// 		if password != "" {
// 			_, err = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", newDevice.Name, "password", password).CombinedOutput()
// 		} else {
// 			_, err = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", newDevice.Name).CombinedOutput()
// 		}
// 		if ctx.Err() == context.DeadlineExceeded {
// 			return connectDeviceMsg{err: fmt.Errorf("Error: connection timed out")}
// 		}
// 		if err != nil {
// 			return connectDeviceMsg{err: err}
// 		}
// 		return connectDeviceMsg{output: "connection successfully activated"}
// 	}
// }
//
//
// // UI
