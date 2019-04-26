package main

import(
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"./trace"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))
	})

	data := map[string]interface{}{
		"Host":r.Host,
	}
	if authCookie,err := r.Cookie("auth");err == nil{
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w,data)
}


func main(){
	var addr = flag.String("host",":8080","localhost")
	flag.Parse()
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		google.New("291466367343-fm01nneeknrqdrmomot7r11rtaq313p1.apps.googleusercontent.com","z0c3sYBv6-eu3hpr7JokPywL","http://localhost:8080/auth/callback/google"),
		)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat",MustAuth(&templateHandler{filename:"chat.html"}))
	http.Handle("/login",&templateHandler{filename:"login.html"})
	http.HandleFunc("/auth/",loginHandler)
	http.Handle("/room",r)
	//チャットルームを開始
	go r.run()
	//Webサーバを起動
	log.Println("Webサーバーを開始します．ポート: ",*addr)

	if err := http.ListenAndServe(*addr,nil);err != nil{
		log.Fatal("ListenAndServe",err)
	}
}