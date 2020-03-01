---
title: "{{.Article.Title}}"
date: {{.Article.PublishedOn.Format "02 Jan 06 15:04 MST"}}
{{if .Article.HeaderImageURL}}
featured_image: "{{.Article.HeaderImageURL}}"
{{end}}
draft: false
---

{{.Article.HTMLBody}}