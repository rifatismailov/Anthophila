package information

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type Info struct {
}

func NewInfo() *Info {
	return &Info{}
}

func (i Info) InfoJson() string {

	type message struct {
		HostName    string `json:"HostName"`
		HostAddress string `json:"HostAddress"`
		MACAddress  string `json:"MACAddress"`
		RemoteAddr  string `json:"RemoteAddr"`
	}

	// Заповнення внутрішньої структури
	msg := message{
		HostName:    i.HostName(),
		HostAddress: i.HostAddress(),
		MACAddress:  i.GetMACAddress(),
		RemoteAddr:  i.RemoteAddress("https://api.ipify.org"),
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
	}
	return string(jsonData)
}

func (i Info) HostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Помилка отримання хост-імені"
	}
	return hostname

}
func (i Info) HostAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Помилка отримання локальних IP-адрес"
	}

	for _, addr := range addrs {
		// Перевіряємо, чи це IP-адреса, а не MAC або інша адреса
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String() // Повертаємо локальну IP-адресу
			}
		}
	}
	return "IP-адреса не знайдена"
}
func (i Info) GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "помилка отримання мережевих інтерфейсів: %v"
	}

	for _, interf := range interfaces {
		// Перевіряємо, чи інтерфейс має MAC-адресу
		if len(interf.HardwareAddr) > 0 {
			return interf.HardwareAddr.String()
		}
	}
	return "MAC-адреса не знайдена"
}
func (i Info) RemoteAddress(urlSite string) string {
	resp, err := http.Get(urlSite)
	if err != nil {
		return "Помилка отримання зовнішньої IP-адреси"
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Помилка читання відповіді"
	}
	// повертаємо Зовнішня IP-адреса:", string(ip)
	return string(ip)
}
