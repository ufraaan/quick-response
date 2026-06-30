package main

import (
	"encoding/json"
	"image/color"
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
body { font-family: -apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif; display: flex; min-height: 100vh; background: #000; color: #fff; }
.left { flex: 1; position: relative; display: flex; justify-content: center; align-items: center; background: #050505; overflow: hidden; }
.hero { position: relative; z-index: 2; width: 75%; max-width: 500px; aspect-ratio: 1; display: flex; justify-content: center; align-items: center; }
.hero svg { width: 100%; height: 100%; display: block; overflow: visible; }
.hero svg polygon { cursor: pointer; }
#glow-g { transform-origin: 0 0; }
.noise { position: absolute; inset: 0; z-index: 1; pointer-events: none; background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.06'/%3E%3C/svg%3E"); background-size: 96px 96px; }
.right { flex: 1; display: flex; justify-content: center; align-items: center; background: #050505; }
@media (max-width: 768px) { body { flex-direction: column; } .left { flex: none; height: 50vh; } .hero { width: 60%; } }
.card { width: 100%; max-width: 340px; padding: 1.5rem; }
h1 { font-size: 1.25rem; font-weight: 500; margin-bottom: 1.5rem; letter-spacing: -0.02em; }
label { display: block; font-size: 0.8125rem; color: #888; margin-bottom: 0.375rem; }
input, select { width: 100%; padding: 0.5rem 0.75rem; margin-bottom: 1rem; background: #111; border: 1px solid #222; border-radius: 6px; color: #fff; font-size: 0.875rem; font-family: inherit; outline: none; }
input:focus, select:focus { border-color: #555; }
input::placeholder { color: #555; }
select { cursor: pointer; }
select option { background: #111; }
button { width: 100%; padding: 0.5rem 0.75rem; background: #fff; color: #000; border: none; border-radius: 6px; font-size: 0.875rem; font-weight: 500; cursor: pointer; font-family: inherit; }
button:hover { background: #e5e5e5; }
button:active { background: #ccc; }
#result { margin-top: 1.5rem; text-align: center; min-height: 4rem; }
#result img { max-width: 100%; border-radius: 4px; }
.footer { margin-top: 3rem; text-align: center; font-size: 0.75rem; color: #666; }
.footer a { color: #60a5fa; text-decoration: underline; text-decoration-color: #1e40af; text-underline-offset: 3px; }
.footer a:hover { color: #93c5fd; text-decoration-color: #3b82f6; }
</style>
</head>
<body>
<div class="left">
<div class="noise"></div>
<div class="hero">
<svg viewBox="-150 -150 300 300" xmlns="http://www.w3.org/2000/svg" style="overflow:visible">
<defs>
<filter id="g1" filterUnits="userSpaceOnUse" x="-150" y="-150" width="300" height="300"><feGaussianBlur in="SourceGraphic" stdDeviation="4"/></filter>
<filter id="g2" filterUnits="userSpaceOnUse" x="-150" y="-150" width="300" height="300"><feGaussianBlur in="SourceGraphic" stdDeviation="12"/></filter>
<filter id="g3" filterUnits="userSpaceOnUse" x="-200" y="-200" width="400" height="400"><feGaussianBlur in="SourceGraphic" stdDeviation="25"/></filter>
<filter id="g4" filterUnits="userSpaceOnUse" x="-250" y="-250" width="500" height="500"><feGaussianBlur in="SourceGraphic" stdDeviation="45"/></filter>
<filter id="g5" filterUnits="userSpaceOnUse" x="-300" y="-300" width="600" height="600"><feGaussianBlur in="SourceGraphic" stdDeviation="70"/></filter>
</defs>
<g id="glow-g">
<polygon points="0,-78 68,39 -68,39" fill="#ffffff" filter="url(#g5)" opacity="0.03"/>
<polygon points="0,-78 68,39 -68,39" fill="#ffffff" filter="url(#g4)" opacity="0.06"/>
<polygon points="0,-78 68,39 -68,39" fill="#ffffff" filter="url(#g3)" opacity="0.12"/>
<polygon points="0,-78 68,39 -68,39" fill="#ffffff" filter="url(#g2)" opacity="0.3"/>
<polygon points="0,-78 68,39 -68,39" fill="#ffffff" filter="url(#g1)" opacity="0.65"/>
</g>
<polygon points="0,-78 68,39 -68,39" fill="#000000"/>
</svg>
</div>
</div>
<div class="right">
<div class="card">
<h1>qr code generator</h1>
<label for="txt">text</label>
<input type="text" id="txt" placeholder="enter text or url" value="hello world">
<label for="sz">size</label>
<select id="sz">
<option value="128">128</option>
<option value="256" selected>256</option>
<option value="512">512</option>
<option value="1024">1024</option>
</select>
<button onclick="gen()">generate</button>
<div id="result"></div>
<div class="footer">100% Go, running on vercel with <a href="https://vercel.com/blog/dockerfile-on-vercel">dockerfile.vercel ↗</a></div>
</div>
</div>
<script>
var gg=document.getElementById('glow-g'),left=document.querySelector('.left');
var sx=0,sy=0,tx=0,ty=0,ms=5;
var per=406.6,pos=0,sc=1,sct=1,flash=0,op=1;
var edges=[
{x1:0,y1:-78,x2:68,y2:39,len:135.3},
{x1:68,y1:39,x2:-68,y2:39,len:136},
{x1:-68,y1:39,x2:0,y2:-78,len:135.3}
];
var onorms=[[0.756,0.654],[0,-1],[-0.756,0.654]];
function gp(pp){
for(var i=0;i<3;i++){var e=edges[i];if(pp<=e.len){var t=pp/e.len;return{x:e.x1+(e.x2-e.x1)*t,y:e.y1+(e.y2-e.y1)*t,ei:i};}pp-=e.len;}
return{x:0,y:-78,ei:0};}
var vx=[0,68,-68],vy=[-78,39,39];
function pit(px,py){
function s(px1,py1,i1,i2){return(px1-vx[i2])*(vy[i1]-vy[i2])-(vx[i1]-vx[i2])*(py1-vy[i2]);}
var d1=s(px,py,0,1),d2=s(px,py,1,2),d3=s(px,py,2,0);
return !((d1<0||d2<0||d3<0)&&(d1>0||d2>0||d3>0));}
var cn=false,ca=0,cd=0,md,ov=false,sv=left.querySelector('svg');

var ns=24,rcircs=[],sn='http://www.w3.org/2000/svg';
for(var i=0;i<ns;i++){
var t=i/ns,pp=gp(t*per),no=onorms[pp.ei];
var cx=pp.x+no[0]*5,cy=pp.y+no[1]*5;
var hue=Math.round((t*360)%360);
var c=document.createElementNS(sn,'circle');
c.setAttribute('cx',cx);c.setAttribute('cy',cy);
c.setAttribute('r',8);
c.setAttribute('fill','hsl('+hue+',100%,60%)');
c.setAttribute('filter','url(#g2)');
c.setAttribute('opacity',0);
c.style.pointerEvents='none';
gg.appendChild(c);rcircs.push(c);}
var rbowA=false,rbowT=0;

onmousemove=function(e){
var r=left.getBoundingClientRect();md=r.width*0.35;
var dx=e.clientX-r.left-r.width/2,dy=e.clientY-r.top-r.height/2;
cd=Math.sqrt(dx*dx+dy*dy);ca=Math.atan2(dy,dx);cn=cd<md;
var vs=300/r.width;ov=pit(dx*vs,-dy*vs);
};
left.addEventListener('click',function(e){
var t=e.target;
if(t.tagName=='polygon'&&!t.getAttribute('filter')){flash=1;rbowA=true;rbowT=0;}});
!function anim(){
var nt=Date.now()*0.001;
if(flash>0.01){op=Math.min(1,0.93+flash*0.35);flash*=0.9;}else{flash=0;op=0.93+0.07*Math.sin(nt*0.9);}
gg.setAttribute('opacity',op);
if(rbowA){
rbowT+=0.03;if(rbowT>=1){rbowT=1;rbowA=false;}
var fi=Math.min(1,rbowT*5),fo=Math.max(0,1-(rbowT-0.25)*1.4);
for(var i=0;i<ns;i++)rcircs[i].setAttribute('opacity',Math.min(fi,fo)*0.55);
}else{
for(var i=0;i<ns;i++)rcircs[i].setAttribute('opacity',0);}
sc+=(sct-sc)*0.07;
pos+=1.5;if(pos>per)pos-=per;
var pp=gp(pos),swx=pp.x*0.06,swy=pp.y*0.06;
var t=cn?Math.max(0,1-cd/md):0;t=t*t;
var px=-Math.cos(ca)*t*ms,py=-Math.sin(ca)*t*ms;
if(ov){t=1;px=0;py=0;}
tx=swx*(1-t)+px;ty=swy*(1-t)+py;
sx+=(tx-sx)*0.08;sy+=(ty-sy)*0.08;
gg.setAttribute('transform','translate('+sx.toFixed(2)+','+sy.toFixed(2)+') scale('+sc.toFixed(3)+')');
requestAnimationFrame(anim)}();
</script>
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

	qr, err := qrcode.New(txt, qrcode.Medium)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	qr.ForegroundColor = color.White
	qr.BackgroundColor = color.Transparent

	png, err := qr.PNG(sz)
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
