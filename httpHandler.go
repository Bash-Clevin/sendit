package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func handleLink(c *fiber.Ctx) error {
	link := c.Params("link")
	if _, ok := pipes[link]; !ok {
		return c.Render("link_not_found", link, "layouts/simple")
	}
	link = fmt.Sprintf("/d/%s", link)
	return c.Render("download", link, "layouts/simple_main")
}

func handleDownload(c *fiber.Ctx) error {
	link := c.Params("link")
	pipe, ok := pipes[link]
	if !ok {
		return c.Render("link_not_found", link, "layouts/simple")
	}
	pipe.wch <- c
	c.Set(fiber.HeaderContentDisposition, "attachment; filename")
	c.Set(fiber.HeaderContentType, fiber.MIMEOctetStream)
	<-pipe.donech

	return nil
}

func handleHome(c *fiber.Ctx) error {
	return c.Render("landing/index", nil, "layouts/simple_main")
}
