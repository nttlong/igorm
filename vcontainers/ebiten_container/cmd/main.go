// package main

// import (
// 	// "ebiten_container"
// 	"ebiten_container"
// 	"runtime"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// 	//"fyne.io/fyne/v2/app"
// 	//"fyne.io/fyne/v2/app"
// )

// //& "D:\Program Files\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\VsDevCmd.bat"
// //D:\"Program Files (x86)"\"Microsoft Visual Studio"\2022\BuildTools\Common7\Tools\VsDevCmd.bat

// func main1() {
// 	/*
// 			$env:PATH += ";D:\cygwin64\bin"
// 			D:\cygwin64\bin\gcc.exe
// 			@echo off
// 		call "D:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\VsDevCmd.bat"
// 		set CC=cl
// 		set CXX=cl
// 		set CGO_ENABLED=1
// 		go run main.go

// 	*/
// 	container := ebiten_container.NewEbitenContainer()
// 	container.Run()
// }

/*
$env:PATH += ";D:\mingw64\bin"
$env:CC = "gcc"
$env:CXX = "g++"
$env:CGO_ENABLED = "1"
go run main.go
*/
// func main1() {
// 	fmt.Println("Before window show")
// 	app.New().NewWindow("Test").ShowAndRun()
// 	fmt.Println("After window close")
// }
package main

import (
	_ "ebiten_container"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			myApp.Quit()
		}),
	))

	myWindow.ShowAndRun()
}

/*
$env:PATH += ";D:\mingw64\bin"
$env:CC = "gcc"
$env:CXX = "g++"
$env:CGO_ENABLED = "1"
go clean -cache -modcache -testcache -x
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\go-build"
go build -x -o app.exe main.go
.\app.exe

*/
