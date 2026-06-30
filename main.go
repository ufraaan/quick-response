package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/qr", qr)
	http.HandleFunc("/health", health)

	log.Println("up on :" + port)
	http.ListenAndServe(":"+port, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
<title>quick response</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #f5f5f5; }
.card { background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,0.08); width: 100%; max-width: 400px; }
h1 { font-size: 1.5rem; margin-bottom: 1.5rem; }
input, select { width: 100%; padding: 0.6rem; margin-bottom: 1rem; border: 1px solid #ddd; border-radius: 6px; font-size: 1rem; }
button { width: 100%; padding: 0.6rem; background: #000; color: white; border: none; border-radius: 6px; font-size: 1rem; cursor: pointer; }
button:hover { opacity: 0.8; }
#result { margin-top: 1.5rem; text-align: center; }
#result img { max-width: 100%; }
.label { font-size: 0.875rem; color: #666; margin-bottom: 0.25rem; }
</style>
</head>
<body>
<div class="card">
<h1>qr code generator</h1>
<p class="label">text</p>
<input type="text" id="txt" placeholder="enter text or url" value="hello world">
<p class="label">size</p>
<select id="sz">
<option value="128">128</option>
<option value="256" selected>256</option>
<option value="512">512</option>
<option value="1024">1024</option>
</select>
<button onclick="gen()">generate</button>
<div id="result"></div>
</div>
<script>
function gen() {
var txt = document.getElementById('txt').value
var sz = document.getElementById('sz').value
if (!txt) return
var img = document.createElement('img')
img.src = '/qr?text=' + encodeURIComponent(txt) + '&size=' + sz
img.alt = 'qr for ' + txt
document.getElementById('result').innerHTML = ''
document.getElementById('result').appendChild(img)
}
</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func qr(w http.ResponseWriter, r *http.Request) {
	txt := r.URL.Query().Get("text")
	if txt == "" {
		http.Error(w, "text param required", http.StatusBadRequest)
		return
	}

	sz := 256
	if s := r.URL.Query().Get("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 && n <= 2048 {
			sz = n
		}
	}

	png, err := qrcode.Encode(txt, qrcode.Medium, sz)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	w.Write(png)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
