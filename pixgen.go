package main

import (
	"fmt"
	"image"
	_"image/jpeg"
	_"image/png"
	"os"
	"github.com/disintegration/imaging"
	"strconv"
	"github.com/fatih/color"
)


func main() {
	yellsty := color.New(color.FgYellow, color.Bold)
	yellsty.Println(` 
	 ______   ______   ______   __   __       ______   ______  ______  
	/\  __ \ /\  ___\ /\  ___\ /\ \ /\ \     /\  __ \ /\  == \/\__  _\ 
	\ \  __ \\ \___  \\ \ \____\ \ \\ \ \    \ \  __ \\ \  __<\/_/\ \/ 
	 \ \_\ \_\\/\_____\\ \_____\\ \_\\ \_\    \ \_\ \_\\ \_\ \_\ \ \_\ 
	  \/_/\/_/ \/_____/ \/_____/ \/_/ \/_/     \/_/\/_/ \/_/ /_/  \/_/ 
																	   `)
	fmt.Print("Image path: ")
    var path string
	fmt.Scanln(&path)
	
	//open image
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}
	defer f.Close()

	//decode image
	im, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}

	//resize if needed
	fmt.Print("Resize mul: ")
    var sizemulstr string
	fmt.Scanln(&sizemulstr) 
	size := im.Bounds().Size()
	sizemul, err := strconv.ParseFloat(sizemulstr, 64)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}
	resize := []int{int(float64(size.X)*sizemul), int(float64(size.Y)*sizemul)}
	reim := imaging.Resize(im, resize[0], resize[1], imaging.Lanczos)
	reimSize := reim.Bounds().Size()

	//setup
	tone := []int{13106, 26213, 39319, 52425, 65534}
	//strTone := []string{"██","▓▓","▒▒","░░","  "}
	strTone := []string{"&&","$$","MM","OO","__"}
	strArt := ""
	fmt.Println("reading and converting pixel...")
	//read RGBA value pixel by pixel => convert to gray value => add to string 
	for y := 0; y < reimSize.Y; y++ {
		for x := 0; x < reimSize.X; x++ {
			fmt.Print("\ry:", y+1, "/", reimSize.Y, "  x:", x+1, "/", reimSize.X)
			r, g, b, _ := reim.At(x,y).RGBA()
			gray := 0.299 * float64(r) +  0.587 * float64(g) + 0.114 * float64(b)
			grayInt := int(gray)
			
			if grayInt <= tone[0] {
				strArt += strTone[0]
			} else if grayInt > tone[0] && grayInt <= tone[1] {
				strArt += strTone[1]
			} else if grayInt > tone[1] && grayInt <= tone[2] {
				strArt += strTone[2]
			} else if grayInt > tone[2] && grayInt <= tone[3] {
				strArt += strTone[3]
			} else {
				strArt += strTone[4]
			}
		}
		strArt += "\n"
	}

	//create .txt file => write string to file
	tex, err := os.Create(f.Name() + ".txt")
	if err!= nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}
	tex.WriteString(strArt)
	//create .html file => write string to file
	html, err := os.Create(f.Name() + ".html")
	if err!= nil {
		fmt.Println("err:", err)
		os.Exit(2)
	}
	html.WriteString(`<button id="conv" style="font-size: 30px">Convert to image</button><pre id="ty" style="display: inline-block">  Powered by <a href="https://html2canvas.hertzen.com/">html2canvas</a></pre>
	<script>
		document.getElementById("conv").addEventListener('click', (e)=>{
		  let temp = document.createElement('pre')
		  temp.textContent = " Loading..."
		  let canvas = html2canvas(document.body)
		  canvas.then((re) => {
			let temp = document.createElement('pre')
			document.body.replaceWith(re)
		  })
		  document.getElementById("conv").replaceWith(temp)
		  document.getElementById("disp").remove()
		  document.getElementById("ty").remove()
		})
	  </script>
	  <script type="text/javascript" src="html2canvas.js"></script>
	  <pre id="disp" style="font-family: Courier; font-size: 10px;">`+"\n" + strArt + "\n</pre>")

	// ***recommend fonts: *Inversionz, courier
	
	fmt.Println("\nComplete...")
	
}	
