package gateway

import "net/http"

func (s *Server) routes() {
	newsRouter := s.Router.Mux.PathPrefix("/api/comments").Subrouter()
	//newsRouter.HandleFunc("/", s.getNewsComments()).Methods("GET")
	//newsRouter.HandleFunc("/", s.createComment()).Methods("POST")
	//newsRouter.HandleFunc("/{id}", s.getSingleComment()).Methods("GET")
	//newsRouter.HandleFunc("/{id}", s.updateComment()).Methods("PATCH")
	//newsRouter.HandleFunc("/{id}", s.deleteComment()).Methods("DELETE")

	//postsRouter.HandleFunc("/{id}/like", s.likePost()).Methods("PATCH")
	//postsRouter.HandleFunc("/{id}/dislike", s.dislikePost()).Methods("PATCH")

	newsRouter.HandleFunc("/{newsuuid}/comments/", s.getNewsComments()).Methods("GET")
	newsRouter.HandleFunc("/{newsuuid}/comments/", s.createComment()).Methods("POST")
	newsRouter.HandleFunc("/{newsuuid}/comments/{id}", s.getSingleComment()).Methods("GET")
	newsRouter.HandleFunc("/{newsuuid}/comments/{id}", s.updateComment()).Methods("PATCH")
	newsRouter.HandleFunc("/{newsuuid}/comments/{id}", s.deleteComment()).Methods("DELETE")

	s.Router.Mux.HandleFunc("/api/user", s.addUser()).Methods("POST")
	s.Router.Mux.HandleFunc("/api/auth/token", s.getUserToken()).Methods("POST")
	s.Router.Mux.HandleFunc("/api/auth/refresh", s.refreshUserToken()).Methods("POST")
	s.Router.Mux.HandleFunc("/api/user/{uid}", s.getUser()).Methods("GET")

	//s.Router.Mux.HandleFunc("/api/oauth/app", s.createApp()).Methods("POST")
	//s.Router.Mux.HandleFunc("/api/oauth/app/{uid}", s.getAppInfo()).Methods("GET")
	//s.Router.Mux.HandleFunc("/api/oauth/authorize", s.getOAuthCode()).Methods("POST")
	//s.Router.Mux.HandleFunc("/api/oauth/token", s.getTokenFromOAuthCode()).Methods("GET")

	s.Router.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Test")) })
}
