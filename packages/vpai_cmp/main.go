package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fasthttp/router"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func serveFile(w io.Writer, setHeader func(key, value string), setStatus func(status int), filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		setStatus(404)
		setHeader("Content-Type", "text/plain")
		w.Write([]byte("File not found"))
		return
	}
	defer f.Close()

	setHeader("Content-Type", "application/pdf")
	setStatus(200)
	io.Copy(w, f)
}
func main1() {
	app := fiber.New()

	// Route để download file (streaming file)
	app.Get("/download/:filename", func(c *fiber.Ctx) error {
		filePath := filepath.Join("./uploads", c.Params("filename"))

		serveFile(
			c, // io.Writer
			func(key, value string) { c.Set(key, value) }, // setHeader
			func(status int) { c.Status(status) },         // setStatus
			filePath,
		)

		return nil
	})

	// Chạy server trên port 8080
	err := app.Listen(":8081")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// file: main_fiber.go

func main2() {
	app := fiber.New()

	app.Get("/download/:filename", func(c *fiber.Ctx) error {
		filePath := filepath.Join("./uploads", c.Params("filename"))
		return c.SendFile(filePath) // fasthttp sendfile
	})

	app.Listen(":8082")
}
func main3() {
	filePath := "./uploads/tes-004.pdf"

	fasthttp.ListenAndServe(":8083", func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/file" {
			f, err := os.Open(filePath)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetBodyString("cannot open file")
				return
			}

			// Lấy thông tin file để set Content-Length
			stat, _ := f.Stat()

			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetContentType("application/pdf")
			ctx.Response.Header.Set("Content-Disposition", "attachment; filename=tes-10kb.pdf")
			ctx.Response.Header.Set("Content-Length", string(stat.Size()))

			// Stream trực tiếp từ file ra response
			ctx.SetBodyStream(f, int(stat.Size()))

			// Không đóng file ở đây, fasthttp sẽ đóng sau khi gửi xong
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	})
}

func main() {
	r := router.New()

	// Route download file PDF
	r.GET("/download/{name}", func(ctx *fasthttp.RequestCtx) {
		name := fmt.Sprintf("%v", ctx.UserValue("name"))
		filePath := filepath.Join("./uploads", name)

		// Kiểm tra tồn tại
		f, err := os.Open(filePath)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("File not found")
			return
		}
		defer f.Close()

		// Header HTTP
		ctx.SetContentType("application/pdf")
		// ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))

		// Copy stream file → response body
		_, err = io.Copy(ctx.Response.BodyWriter(), f)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
	})

	// Start server
	fasthttp.ListenAndServe(":8081", r.Handler)
}
