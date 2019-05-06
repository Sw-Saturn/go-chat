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
	"github.com/joho/godotenv"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

var avatars Avatar = TryAvatars{
	UseGravatar,
	UseFileSystemAvatar,
	UseAuthAvatar,
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env")
	}
	var addr = flag.String("host",":8080","localhost")
	flag.Parse()
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_CLIENT"),os.Getenv("GOOGLE_SECRET"),"http://localhost:8080/auth/callback/google"),
		)
	r := newRoom(avatars)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat",MustAuth(&templateHandler{filename:"chat.html"}))
	http.Handle("/login",&templateHandler{filename:"login.html"})
	http.HandleFunc("/auth/",loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, request *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:"auth",
			Value:"",
			Path:"/",
			MaxAge:-1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room",r)
	http.Handle("/upload",&templateHandler{filename:"upload.html"})
	http.HandleFunc("/uploader",uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	//チャットルームを開始
	go r.run()
	//Webサーバを起動
	log.Println("Webサーバーを開始します．ポート: ",*addr)

	if err := http.ListenAndServe(*addr,nil);err != nil{
		log.Fatal("ListenAndServe",err)
	}
}