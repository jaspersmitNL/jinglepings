package main

import (
	"fmt"
	"image/png"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

func LoadImage(path string) MyImage {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	png, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	myImage := MyImage{}
	myImage.width = png.Bounds().Dx()
	myImage.height = png.Bounds().Dy()

	myImage.pixels = make([]byte, myImage.width*myImage.height*4)

	for x := 0; x < myImage.width; x++ {
		for y := 0; y < myImage.height; y++ {
			r, g, b, a := png.At(x, y).RGBA()
			myImage.pixels[(y*myImage.width+x)*4+0] = byte(r >> 8)
			myImage.pixels[(y*myImage.width+x)*4+1] = byte(g >> 8)
			myImage.pixels[(y*myImage.width+x)*4+2] = byte(b >> 8)
			myImage.pixels[(y*myImage.width+x)*4+3] = byte(a >> 8)
		}
	}

	return myImage
}

func ping(ip string) {

	// Create an ICMP connection
	conn, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Parse the target IPv6 address
	dst, err := net.ResolveIPAddr("ip6", ip)
	if err != nil {
		log.Fatal(err)
	}

	msg := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PING"),
		},
	}
	// Marshal the message
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send the packet
	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		log.Fatalln("Failed to write message:", err)
	}
	// cmd := exec.Command("ping6", "-c", "1", ip)
	// cmd.Run()
}

type PingJob struct {
	ip string
}

type MyImage struct {
	width  int
	height int

	pixels []byte
}

var jobs = make(chan PingJob)
var workers = 2
var wg sync.WaitGroup

// worker function that sends a ping to the ip address
func worker(jobs <-chan PingJob) {
	for job := range jobs {
		ping(job.ip)
	}
	wg.Done()
}

func draw(x, y, r, g, b int) {
	PRE := "2001:610:1908:a000:"
	a := 255

	xHex := fmt.Sprintf("%04x", x)
	yHex := fmt.Sprintf("%04x", y)

	rHex := fmt.Sprintf("%02x", r)
	gHex := fmt.Sprintf("%02x", g)
	bHex := fmt.Sprintf("%02x", b)
	aHex := fmt.Sprintf("%02x", a)

	ip := fmt.Sprintf("%s%s:%s:%s%s:%s%s", PRE, xHex, yHex, bHex, gHex, rHex, aHex)

	job := PingJob{ip}
	jobs <- job
}

func drawImage(image *MyImage, x2, y2 int) {
	for x := 0; x < image.width; x++ {
		for y := 0; y < image.height; y++ {
			r := int(image.pixels[(y*image.width+x)*4+0]) // *4+0
			g := int(image.pixels[(y*image.width+x)*4+1])
			b := int(image.pixels[(y*image.width+x)*4+2])
			draw(x2+x, y2+y, r, g, b)
		}
	}
}

func main() {
	log.Println("Hello, World!")

	// start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs)
	}

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			draw(i, j, 255, 0, 0)
		}
	}

	//image.RGBAd
	// image := LoadImage("/root/rackmeme2.png")

	// for {
	// 	x := rand.IntN(100) // 0-99
	// 	y := rand.IntN(100) // 0-99

	// 	drawImage(&image, x, y)

	// 	time.Sleep(1 * time.Second)
	// }

	close(jobs)
	wg.Wait()

}
