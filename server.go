package emg

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

const maxScores = 20

type HighScore struct {
	Name  string
	Score int
}

type Server struct {
	lock   *sync.Mutex
	scores []HighScore
}

func New() (Server, error) {
	hs, err := deobfScores()
	if err != nil {
		return Server{}, err
	}
	s := Server{
		lock:   &sync.Mutex{},
		scores: hs,
	}
	http.HandleFunc("/scores", s.Scores)
	http.HandleFunc("/", s.HTML)
	return s, nil
}

func (s *Server) HTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "favicon.ico" {
		status := http.StatusNotFound
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Write([]byte(`<html>
	<head>
<script>
const totalQuestions = 20;
let asked = 0;
let correct = 0;

const startGame = function() {
	document.getElementById('start').style.display = "none";
	document.getElementById('question').style.display = "block";
	document.getElementById('final').style.display = "none";
	asked = 0;
	correct = 0;
	nextQuestion();
};

const nextQuestion = function() {
	let o1 = document.getElementById('o1');
	let o2 = document.getElementById('o2');
	let res = document.getElementById('result');
	asked++;
	document.getElementById('qheader').innerText = "Question " + asked + "/" + totalQuestions;
	o1.innerText = Math.floor(Math.random() * 10 + 1);
	o2.innerText = Math.floor(Math.random() * 10 + 1);
	res.value = 0;
}

const answer = function() {
	const o1v = parseInt(document.getElementById('o1').innerText);
	const o2v = parseInt(document.getElementById('o2').innerText);
	const resv = parseInt(document.getElementById('result').value);
	if (o1v * o2v === resv) {
		correct++;
	}
	if (asked < totalQuestions) {
		nextQuestion();
		return;
	}
	document.getElementById('question').style.display = "none";
	document.getElementById('final').style.display = "block";
	document.getElementById('finalScore').innerText = "You got " + correct + "/" + asked;
};

const showScores = function() {
	let xhr = new XMLHttpRequest();
	xhr.open('GET', '/scores');
	xhr.onload = function() {
		const hs = JSON.parse(xhr.response);
		let scores = document.getElementById('scores');
		scores.innerHTML = '<tr><th>Player</th><th>Score</th></tr>';
		for (let i = 0; i < hs.length; i++) {
			scores.innerHTML += "<tr><td>" + hs[i].Name + "</td><td>" + hs[i].Score + "</td></tr>";
		}
	};
	xhr.send();
};

const submit = function() {
	const name = document.getElementById('name').value;
	let xhr = new XMLHttpRequest();
	xhr.open('POST', '/scores');
	xhr.onload = function() {
		const hs = JSON.parse(xhr.response);
		let scores = document.getElementById('scores');
		scores.innerHTML = '<tr><th>Player</th><th>Score</th></tr>';
		for (let i = 0; i < hs.length; i++) {
			scores.innerHTML += "<tr><td>" + hs[i].Name + "</td><td>" + hs[i].Score + "</td></tr>";
		}
		document.getElementById('final').style.display = "none";
		document.getElementById('start').style.display = "block";
	};
	xhr.onerror = function() {
		console.log("error: " + xhr.statusText);
	};
	let formData = new FormData();
	formData.append('name', name);
	formData.append('score', correct);
	console.log("Sending: ", formData);
	xhr.send(formData);
};

</script>
	</head>
	<body onload="showScores()">
		<h1>EPIC MATH GAME</h1>
<div id="start" style="display: block">
	<button type="button" onclick="startGame()">Start Game</button>
</div>
<div id="question" style="display: none">
	<h2 id="qheader"></h2>
	<span id="o1">1</span> x <span id="o2">1</span> = <input id="result" />
	<button type="button" onclick="answer()">Answer</button>
</div>
<div id="final" style="display: none">
	<h2 id="finalScore"></h2>
	<div>Name: <input id="name" /><button type="button" onclick="submit()">Submit Score</button></div>
	<button type="button" onclick="startGame()">Restart</button>
</div>
<div>
	<h2>High Scores</h2>
	<table id="scores">
	</table>
</div>
	</body>
</html>`))
}

func (s *Server) Scores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Print("parsing form: ", err.Error())
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			log.Print("parsing multipart form: ", err.Error())
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		name := r.Form.Get("name")
		scoreStr := r.Form.Get("score")
		score, err := strconv.Atoi(scoreStr)
		if err != nil || name == "" {
			if err != nil {
				log.Print("Bad Request: ", err.Error())
			}
			if name == "" {
				log.Print("empty name")
			}
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		scores := append(s.scores, HighScore{Name: name, Score: score})
		sort.Slice(scores, func(i, j int) bool {
			return !(scores[i].Score < scores[j].Score)
		})
		if len(scores) > maxScores {
			scores = scores[:maxScores]
		}
		s.scores = scores
	}

	b, err := json.Marshal(s.scores)
	if err != nil {
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
