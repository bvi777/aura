package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	data := make([]byte, 0)
	var (
		rate       int
		inByte, cs byte
	)
	i2c := machine.I2C1
	err := i2c.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})
	if err != nil {
		println("i2c error")
	}

	lcd := hd44780i2c.New(machine.I2C1, 0x3f)
	err = lcd.Configure(hd44780i2c.Config{
		Width:       16,
		Height:      2,
		CursorOn:    false,
		CursorBlink: false,
	})
	if err != nil {
		println("lcd error")
	}

	led := machine.NEOPIXEL // for zero
	// led := machine.LED // for pico
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	rgb := ws2812.New(machine.NEOPIXEL)
	color1 := []color.RGBA{color.RGBA{255, 0, 0, 0xff}}
	color2 := []color.RGBA{color.RGBA{0, 255, 0, 0xff}}

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{BaudRate: 9600})

	lcd.Print([]byte("***   AURA   ***\nFlow rate meter"))
	time.Sleep(time.Millisecond * 3000)
	lcd.ClearDisplay()
	lcd.Print([]byte("Air rate:\n0.00 l/min"))

	for {
		// led.Low()
		time.Sleep(time.Millisecond * 100)
		rgb.WriteColors(color1)

		// led.High()
		time.Sleep(time.Millisecond * 100)
		rgb.WriteColors(color2)

		data = nil
		cs = 0
		for uart.Buffered() > 0 {
			inByte, err = uart.ReadByte()
			if err != nil {
				println(err)
			}

			data = append(data, inByte)
			cs -= inByte
		}
		lcd.SetCursor(0, 1)
		if len(data) == 12 && cs == 0 {
			rate = int(data[5])*256 + int(data[6])
			lcd.Print([]byte(fmt.Sprintf("%2.2f l/min  ", float64(rate)/100.0)))

			// } else {
			// lcd.Print([]byte(fmt.Sprintf("%d %d %d ", int(cs), int(data[11]), len(data)))) for debug
		}
		// println(data)
	}
}

// compile & flash:
// tinygo flash -target=waveshare-rp2040-zero
