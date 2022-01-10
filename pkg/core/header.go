package core

import httpdownloader "github.com/linuxsuren/http-downloader/pkg"

// AsFormRequest is the fixed header for the form request purpose
var AsFormRequest = map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}
