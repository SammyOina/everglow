package everglow

import "net"

func getHostIP(host string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

const ESP8266Host = "esp32-arduino.mshome.net"

var UDP_IP = "192.168.137.244"

const UDP_PORT = 7777
const SOFTWARE_GAMMA_CORRECTION = false
const USE_GUI = true
const DISPLAY_FPS = true
const N_PIXELS = 256
const GAMMA_TABLE_PATH = "gamma_table.npy"
const MIC_RATE = 44100
const FPS = 60
const MIN_FREQUENCY = 200
const MAX_FREQUENCY = 12000
const N_FFT_BINS = 24
const N_ROLLING_HISTORY = 2
const MIN_VOLUME_THRESHOLD = 1e-7
