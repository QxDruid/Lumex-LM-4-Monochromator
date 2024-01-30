package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

var COM_PORT string = ""

func go_to(wavelength int) int {

	hexString := fmt.Sprintf("%04X", wavelength)
	hexString = hexString[2:4] + hexString[0:2]
	command := "4C" + hexString + "03\r"

	port, err := serial.Open(COM_PORT, mode)
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
	port, err := serial.Open(COM_PORT, mode)
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

func convert_state_to_wavelength(state string) int {
	state_split := strings.Fields(state)
	if len(state_split) < 3 {
		return 0
	}
	wavelength := state_split[2] + state_split[1]
	wavelength_int, err := strconv.ParseInt(wavelength, 16, 32)
	if err != nil {
		return 0
	}

	return int(wavelength_int)
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

func input_wl_to_int(input string) int {
	err_ := input_wl_validator(input)
	if err_ == 0 {

		return 0
	}
	wavelength_float, err := strconv.ParseFloat(input, 8)
	if err != nil {
		log.Fatal(err)
	}

	wavelength_float *= 10
	return int(wavelength_float)
}

func main() {
	a := app.New()
	window := a.NewWindow("test")
	window.Resize(fyne.NewSize(320, 520))
	label_status := widget.NewLabel("Is OK!")
	label_status.Move(fyne.NewPos(180, 460))
	label_status.Resize(fyne.NewSize(140, 40))

	// Get avaliable COM
	menu_com := fyne.NewMenu("COM Port")
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		label_status.SetText("No serial found!")
	} else {
		for _, port := range ports {
			actions_item := fyne.NewMenuItem(port, func() {
				COM_PORT = port
				text := fmt.Sprintf("%v", COM_PORT)
				if err != nil {
					log.Fatal(err)
				}
				label_status.SetText(text)
			})
			menu_com.Items = append(menu_com.Items, actions_item)
		}
	}

	main_menu := fyne.NewMainMenu(menu_com)
	window.SetMainMenu(main_menu)

	label_1 := widget.NewLabel("Current WL (nm):")
	label_1.Move(fyne.NewPos(60, 20))
	label_1.Resize(fyne.NewSize(140, 40))

	label_current_wl := widget.NewLabel("400.0")
	label_current_wl.Move(fyne.NewPos(180, 20))
	label_current_wl.Resize(fyne.NewSize(100, 40))

	label_go_wl := widget.NewLabel("Go to WL:")
	label_go_wl.Move(fyne.NewPos(40, 60))
	label_go_wl.Resize(fyne.NewSize(100, 40))

	entry_go_wl := widget.NewEntry()
	entry_go_wl.Text = "400.0"
	entry_go_wl.Move(fyne.NewPos(160, 60))
	entry_go_wl.Resize(fyne.NewSize(100, 40))

	btn_go_to_wl := widget.NewButton("Go", func() {
		if COM_PORT == "" {
			label_status.SetText("SELECT COM")
			return
		}

		data := entry_go_wl.Text
		wavelength_to_go := input_wl_to_int(data)
		if wavelength_to_go == 0 {

			label_status.SetText("False Wavelength")
			return
		}
		var current_wl int
		for current_wl == 0 {

			current_wl = convert_state_to_wavelength(read_state())
		}
		if wavelength_to_go == current_wl {
			label_status.SetText("Is OK!")
			return
		}

		label_status.SetText("In Process")
		go_to(wavelength_to_go)

		for current_wl != wavelength_to_go {
			current_wl = convert_state_to_wavelength(read_state())
			if current_wl == 0 {
				continue
			}
			current_wl_str := strconv.Itoa(current_wl)
			current_wl_str = current_wl_str[:3] + string('.') + string(current_wl_str[3])

			label_current_wl.SetText(current_wl_str)
		}
		label_status.SetText("Is OK!")

	})
	btn_go_to_wl.Move(fyne.NewPos(40, 120))
	btn_go_to_wl.Resize(fyne.NewSize(100, 40))

	btn_read_state := widget.NewButton("Get", func() {
		if COM_PORT == "" {
			label_status.SetText("SELECT COM")
			return
		}

		data := convert_state_to_wavelength(read_state())
		s2 := strconv.Itoa(data)
		s2 = s2[:3] + string('.') + string(s2[3])
		label_current_wl.SetText(s2)
	})
	btn_read_state.Move(fyne.NewPos(160, 120))
	btn_read_state.Resize(fyne.NewSize(100, 40))

	label_go_from_wl := widget.NewLabel("Go From (nm):")
	label_go_from_wl.Move(fyne.NewPos(40, 180))
	label_go_from_wl.Resize(fyne.NewSize(60, 20))

	entry_go_from_wl := widget.NewEntry()
	entry_go_from_wl.Text = "400.0"
	entry_go_from_wl.Move(fyne.NewPos(160, 180))
	entry_go_from_wl.Resize(fyne.NewSize(100, 40))

	label_go_to_wl := widget.NewLabel("Go To (nm):")
	label_go_to_wl.Move(fyne.NewPos(40, 240))
	label_go_to_wl.Resize(fyne.NewSize(60, 20))

	entry_go_to_wl := widget.NewEntry()
	entry_go_to_wl.Text = "400.0"
	entry_go_to_wl.Move(fyne.NewPos(160, 240))
	entry_go_to_wl.Resize(fyne.NewSize(100, 40))

	label_delay := widget.NewLabel("Step Delay (ms)")
	label_delay.Move(fyne.NewPos(40, 300))
	label_delay.Resize(fyne.NewSize(60, 20))

	entry_delay := widget.NewEntry()
	entry_delay.Text = "100"
	entry_delay.Move(fyne.NewPos(160, 300))
	entry_delay.Resize(fyne.NewSize(100, 40))

	label_step := widget.NewLabel("Step (nm)")
	label_step.Move(fyne.NewPos(40, 360))
	label_step.Resize(fyne.NewSize(60, 20))

	select_step := widget.NewSelect(
		[]string{
			"10",
			"1",
			"0.5",
			"0.1",
		},
		func(s string) {

		},
	)
	select_step.PlaceHolder = "0.1"
	select_step.Selected = "0.1"
	select_step.Move(fyne.NewPos(160, 360))
	select_step.Resize(fyne.NewSize(100, 40))

	btn_go_from_to := widget.NewButton("Go", func() {

		var current_wl int = 0

		if COM_PORT == "" {
			label_status.SetText("SELECT COM")
			return
		}

		// Read and parse forms
		// From wavelength value
		data_from := input_wl_to_int(entry_go_from_wl.Text)
		if data_from == 0 {
			label_status.SetText("False From WL ")
			return
		}

		// To wavelength value
		data_to := input_wl_to_int(entry_go_to_wl.Text)
		if data_to == 0 {
			label_status.SetText("False To WL ")
			return
		}
		// delay ms value
		data_delay := entry_delay.Text
		delay_ms_int, err := strconv.ParseInt(data_delay, 10, 16)
		if err != nil {
			log.Fatal(err)
		}

		// Wavelength Step value
		data_step := select_step.Selected
		data_step_float, err := strconv.ParseFloat(data_step, 8)
		if err != nil {
			log.Fatal(err)
		}
		data_step_int := int(data_step_float * 10)

		label_status.SetText("Go to Start")

		// Go to FROM position
		go_to(data_from)
		current_wl = convert_state_to_wavelength(read_state())
		for current_wl != data_from {
			current_wl = convert_state_to_wavelength(read_state())
			if current_wl == 0 {
				continue
			}
			current_wl_str := strconv.Itoa(current_wl)
			current_wl_str = current_wl_str[:3] + string('.') + string(current_wl_str[3])
			label_current_wl.SetText(current_wl_str)
		}

		label_status.SetText("In Process!")

		// Reverse step
		if data_from > data_to {
			data_step_int *= -1
		}

		// Go to end position with some step
		var next_step int = data_from
		for next_step != data_to {
			// Читаем текущую длину волны
			for current_wl == 0 {
				current_wl = convert_state_to_wavelength(read_state())
			}
			current_wl_str := strconv.Itoa(current_wl)
			current_wl_str = current_wl_str[:3] + string('.') + string(current_wl_str[3])

			label_current_wl.SetText(current_wl_str)

			// Стоим на текущей длине волны нудное время
			time.Sleep(time.Duration(delay_ms_int) * time.Millisecond)

			// Едем на следующую длину волны
			next_step += data_step_int
			go_to(next_step)

			// ждем пока доедем до длины волны шага
			for current_wl != next_step {
				current_wl = convert_state_to_wavelength(read_state())
				for current_wl == 0 {
					current_wl = convert_state_to_wavelength(read_state())
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
		for current_wl == 0 {

			current_wl = convert_state_to_wavelength(read_state())
		}
		current_wl_str := strconv.Itoa(current_wl)
		current_wl_str = current_wl_str[:3] + string('.') + string(current_wl_str[3])
		label_current_wl.SetText(current_wl_str)
		label_status.SetText("Is OK!")

	})

	btn_go_from_to.Move(fyne.NewPos(40, 420))
	btn_go_from_to.Resize(fyne.NewSize(240, 40))

	window_content := container.NewWithoutLayout(label_1, label_current_wl, label_go_wl, entry_go_wl, btn_go_to_wl, btn_read_state, label_go_from_wl,
		entry_go_from_wl, label_go_to_wl, entry_go_to_wl, btn_go_from_to, label_delay, entry_delay, label_status, label_step, select_step)
	window.SetContent(window_content)
	window.ShowAndRun()
}
