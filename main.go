package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

var mode = &serial.Mode{
	BaudRate: 19200,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.OneStopBit,
}

func go_to(wavelength int) int {

	hexString := fmt.Sprintf("%04X", wavelength)
	hexString = hexString[2:4] + hexString[0:2]
	command := "4C" + hexString + "03\r"

	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
	port.Close()
	return n
}

func read_state() string {
	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
	}

	command := "5201\r"
	n, err := port.Write([]byte(command))
	if err != nil {
		log.Fatal(err, n)
	}
	responce := make([]byte, 30, 30)
	for {
		buff := make([]byte, 30)

		n, err = port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}

		responce = append(responce, buff...)
		if buff[n-1] == '\n' {
			break
		}
	}

	port.Close()
	return string(responce)
}

func convert_state_to_wavelength(state string) string {
	state_split := strings.Fields(state)
	wavelength := state_split[2] + state_split[1]
	wavelength_int, err := strconv.ParseInt(wavelength, 16, 32)
	wavelength_int_1 := wavelength_int / 10
	wavelength_int_2 := wavelength_int % wavelength_int_1
	if err != nil {
		panic(err)
	}
	fmt.Print(wavelength_int_1)
	fmt.Print(wavelength_int_2)
	wavelength_str := fmt.Sprintf("%d.%d", wavelength_int_1, wavelength_int_2)
	return wavelength_str
}

func input_wl_validator(input string) int {

	input_float, err := strconv.ParseFloat(input, 8)
	if err != nil {
		return 0
	}
	input_int := int(input_float * 10)
	if input_int > 9000 || input_int < 2200 {
		return 0
	}
	fmt.Println()
	return 1
}

func main() {
	a := app.New()
	window := a.NewWindow("test")
	window.Resize(fyne.NewSize(320, 450))

	label_status := widget.NewLabel("Is OK!")
	label_status.Move(fyne.NewPos(40, 400))
	label_status.Resize(fyne.NewSize(140, 40))

	label_1 := widget.NewLabel("Current WL (nm):")
	label_1.Move(fyne.NewPos(60, 20))
	label_1.Resize(fyne.NewSize(140, 40))
	label_current_wl := widget.NewLabel("400")
	label_current_wl.Move(fyne.NewPos(180, 20))
	label_current_wl.Resize(fyne.NewSize(100, 40))

	label_go_wl := widget.NewLabel("Go to WL:")
	label_go_wl.Move(fyne.NewPos(40, 60))
	label_go_wl.Resize(fyne.NewSize(100, 40))

	entry_go_wl := widget.NewEntry()
	entry_go_wl.Move(fyne.NewPos(160, 60))
	entry_go_wl.Resize(fyne.NewSize(100, 40))

	btn_go_to_wl := widget.NewButton("Go", func() {
		data := entry_go_wl.Text
		err_ := input_wl_validator(data)
		if err_ == 0 {
			label_status.SetText("False Wavelength")
			return
		}
		wavelength_to_go, err := strconv.ParseFloat(data, 8)
		if err != nil {
			log.Fatal(err)
		}
		label_status.SetText("In Process")
		wavelength_to_go *= 10
		go_to(int(wavelength_to_go))
		current_wl := read_state()
		current_wl = convert_state_to_wavelength(data)
		label_current_wl.SetText(current_wl)
		label_status.SetText("Is Ok")

	})
	btn_go_to_wl.Move(fyne.NewPos(40, 120))
	btn_go_to_wl.Resize(fyne.NewSize(100, 40))

	btn_read_state := widget.NewButton("Get", func() {
		data := read_state()
		data = convert_state_to_wavelength(data)
		label_current_wl.SetText(data)
	})
	btn_read_state.Move(fyne.NewPos(160, 120))
	btn_read_state.Resize(fyne.NewSize(100, 40))

	label_go_from_wl := widget.NewLabel("Go From:")
	label_go_from_wl.Move(fyne.NewPos(40, 180))
	label_go_from_wl.Resize(fyne.NewSize(60, 20))

	entry_go_from_wl := widget.NewEntry()
	entry_go_from_wl.Move(fyne.NewPos(160, 180))
	entry_go_from_wl.Resize(fyne.NewSize(100, 40))

	label_go_to_wl := widget.NewLabel("Go To:")
	label_go_to_wl.Move(fyne.NewPos(40, 240))
	label_go_to_wl.Resize(fyne.NewSize(60, 20))

	entry_go_to_wl := widget.NewEntry()
	entry_go_to_wl.Move(fyne.NewPos(160, 240))
	entry_go_to_wl.Resize(fyne.NewSize(100, 40))

	btn_go_from_to := widget.NewButton("Go", func() {
		data_from := entry_go_from_wl.Text
		err_ := input_wl_validator(data_from)
		if err_ == 0 {
			label_status.SetText("False Wavelength From Go")
			return
		}

		data_to := entry_go_to_wl.Text
		err_ = input_wl_validator(data_to)
		if err_ == 0 {
			label_status.SetText("False Wavelength To Go")
			return
		}
		label_status.SetText("In Process")

		//go_from_to(from_int, to_int)
		label_status.SetText("Is Ok")
	})

	label_delay := widget.NewLabel("Step Delay")
	label_delay.Move(fyne.NewPos(40, 300))
	label_delay.Resize(fyne.NewSize(60, 20))

	entry_delay := widget.NewEntry()
	entry_delay.Move(fyne.NewPos(160, 300))
	entry_delay.Resize(fyne.NewSize(100, 40))

	btn_go_from_to.Move(fyne.NewPos(40, 360))
	btn_go_from_to.Resize(fyne.NewSize(240, 40))

	window_content := container.NewWithoutLayout(label_1, label_current_wl, label_go_wl, entry_go_wl, btn_go_to_wl, btn_read_state, label_go_from_wl,
		entry_go_from_wl, label_go_to_wl, entry_go_to_wl, btn_go_from_to, label_delay, entry_delay, label_status)
	window.SetContent(window_content)
	window.ShowAndRun()
}
