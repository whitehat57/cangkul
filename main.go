package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var proxyURLs = []string{
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
	"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/http.txt",
	"https://raw.githubusercontent.com/officialputuid/KangProxy/KangProxy/http/http.txt",
	"https://raw.githubusercontent.com/prxchk/proxy-list/main/http.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/proxylist-to/proxy-list/main/http.txt",
	"https://raw.githubusercontent.com/yuceltoluyag/GoodProxy/main/raw.txt",
	"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/http.txt",
	"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/https.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/https.txt",
	"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/http_proxies.txt",
	"https://raw.githubusercontent.com/opsxcq/proxy-list/master/list.txt",
	"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/https_proxies.txt",
	"https://api.openproxylist.xyz/http.txt",
	"https://api.proxyscrape.com/v2/?request=displayproxies",
	"https://api.proxyscrape.com/?request=displayproxies&proxytype=http",
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all",
	"https://www.proxydocker.com/en/proxylist/download?email=noshare&country=all&city=all&port=all&type=all&anonymity=all&state=all&need=all",
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=anonymous",
	"http://worm.rip/http.txt",
	"https://proxyspace.pro/http.txt",
	"https://multiproxy.org/txt_all/proxy.txt",
	"https://proxy-spider.com/api/proxies.example.txt",
}

func main() {
	fileName := "proxy.txt"

	// Hapus file lama jika ada.
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("File %s sudah ada! Menghapus dan mengunduh yang baru.\n", fileName)
		if err := os.Remove(fileName); err != nil {
			fmt.Printf("Gagal menghapus file: %v\n", err)
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup
	ch := make(chan string) // Channel untuk komunikasi antar goroutine.

	// Membuka file untuk menulis proxy.
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Gagal membuat file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Goroutine untuk menulis proxy ke file.
	go func() {
		for proxy := range ch {
			_, err := file.WriteString(proxy)
			if err != nil {
				fmt.Printf("Gagal menulis ke file: %v\n", err)
			}
		}
	}()

	// Ambil setiap URL secara paralel.
	for _, url := range proxyURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := fetchProxy(url, ch); err != nil {
				fmt.Printf("Gagal mengambil %s: %v\n", url, err)
			}
		}(url)
	}

	wg.Wait()
	close(ch)

	// Hitung total proxy yang diunduh.
	total := countLines(fileName)
	fmt.Printf("\n(%d) Proxy berhasil diunduh.\n", total)
}

// fetchProxy mengunduh proxy dari URL dan mengirimkannya ke channel.
func fetchProxy(url string, ch chan<- string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ch <- string(body)
	fmt.Printf("Mengambil %s\n", url)
	return nil
}

// countLines menghitung jumlah baris dalam file.
func countLines(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Gagal membuka file: %v\n", err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}
