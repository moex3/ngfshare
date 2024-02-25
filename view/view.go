package view

import (
    "log"
    "fmt"
    "io"
    "html/template"
    "time"

    "git.fuwafuwa.moe/x3/ngfshare/config"
)

var tmpl *template.Template

func addFuncs(t *template.Template) *template.Template {
    fmap := template.FuncMap{
        "formatDate": func(unix int64) string {
            return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
        },
        "formatFileSize": func(b int64) string {
            const unit = 1024
            if b < unit {
                return fmt.Sprintf("%d B", b)
            }
            div, exp := int64(unit), 0
            for n := b / unit; n >= unit; n /= unit {
                div *= unit
                exp++
            }
            return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
        },
    }
    return t.Funcs(fmap)
}

func LoadTemplates() error {
    t, err := addFuncs(template.New("")).ParseGlob(fmt.Sprintf("%s/*", config.Conf.HTMLTemplateDir))
    if err != nil {
        log.Println("LoadTemplates:", err)
        return err
    }
    tmpl = t

    log.Println("Templates loaded")

    return nil
}

func Execute(tstr string, wr io.Writer, data any) error {
    return tmpl.ExecuteTemplate(wr, tstr, data)
}
