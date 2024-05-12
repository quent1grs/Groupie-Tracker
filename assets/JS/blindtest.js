var ws = new WebSocket('ws://localhost:8080/chatblindtestws')
        var nom = prompt("Entrez votre nom : ");
        
        ws.onopen = () => {
            console.log("Connected")
            document.getElementById("sendtextchat").addEventListener("click", (e)=> {
                 e.preventDefault();
                 var data = {
                    "message" : document.getElementById("textchat").value,
                }
                ws.send(JSON.stringify(data));
            })
        }

        ws.onmessage = (e) => {
            document.getElementById("chat").innerHTML += e.data + "<br>";
        }

        ws.onerror = (e) => {
            console.log(e)
        }

        ws.onclose = () => {
            delete ws;
        }

        var musicCountdown = 30;
        
        window.onload = function() {
            var conn = new WebSocket('ws://localhost:8080/blindtestws');
            console.log(conn);
    
            var firstload = true;
            var firstMessage = true;

            conn.onopen = function(e) {
                console.log("Connection established!");
            };

            document.getElementById('play').addEventListener('click', function() {
                conn.send('start_music');
            });

            conn.onclose = function(e) {
                console.log("Connection closed!");
            };

            var title;
            var artist;

            conn.onmessage = function(e) {
                if (firstMessage === true) {
                    firstMessage = false;
                    return;
                }
                var data;
                var previousSong;
                var countdownElement = document.getElementById('countdown');

                try {
                    data = JSON.parse(e.data);
                } catch (error) {
                    data = e.data;
                }
                console.log(data.Artist);
                console.log(data.Title);
                console.log(title);
                console.log(artist);

                var audio = document.getElementById('audio');
                var countdownElement = document.getElementById('countdown');

                if (data === 'start_music') {
                    var countdown = 3;
                    var musicCountdownInterval;
                    countdownElement.innerText = countdown;

                    var countdownInterval = setInterval(function() {
                        countdown--;
                        countdownElement.innerText = countdown;

                        if (countdown <= 0) {
                            clearInterval(countdownInterval);
                            audio.play();
                            musicCountdown = 30;
                            countdownElement.innerText = musicCountdown;

                            musicCountdownInterval = setInterval(function() {
                                musicCountdown--;
                                countdownElement.innerText = musicCountdown;

                                if (musicCountdown <= 0) {
                                    clearInterval(musicCountdownInterval);
                                    var previous = document.getElementById('previous');
                                    previous.innerText = "Previous song: " + title + " by " + artist;
                                    conn.send('end_music');
                                }
                            }, 1000);
                        }
                    }, 1000);
                } else if (data.Preview) {
                    title = data.Title;
                    artist = data.Artist;
                    if (firstload === true) {
                        audio.src = data.Preview;
                        audio.load();
                        firstload = false;
                    } else {
                        var countdown = 3;
                        countdownElement.innerText = countdown;
                        
                        var countdownInterval = setInterval(function() {
                            countdown--;
                            countdownElement.innerText = countdown;

                            if (countdown <= 0) {
                                clearInterval(countdownInterval);
                                audio.src = data.Preview;
                                audio.load();
                                musicCountdown = 30;
                                countdownElement.innerText = musicCountdown;

                                musicCountdownInterval = setInterval(function() {
                                    musicCountdown--;
                                    countdownElement.innerText = musicCountdown;

                                    if (musicCountdown <= 0) {
                                        clearInterval(musicCountdownInterval);
                                        var previous = document.getElementById('previous');
                                        previous.innerText = "Previous song: " + title + " by " + artist;
                                        conn.send('end_music');
                                    }
                                }, 1000);
                            }
                        }, 1000);
                    }
                }
            };
    
            document.querySelector('form').addEventListener('submit', function(e) {
                e.preventDefault();
                if (conn.readyState === WebSocket.OPEN) {
                    remainingTime = parseInt(document.getElementById('countdown').innerText, 10);
                    var answer = document.querySelector('input[name="blindtest_answer"]').value;
                    if (typeof answer === 'string') {
                        var message = {
                            answer: answer,
                            remainingTime: remainingTime,
                        };
                        conn.send(JSON.stringify(message));
                    } else {
                        console.error('Answer is not a string:', answer);
                    }
                } else {
                    console.error('WebSocket is not open:', conn.readyState);
                }
            });
        };