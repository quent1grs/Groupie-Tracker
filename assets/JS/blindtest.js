var ws = new WebSocket('ws://localhost:8080/ws')
        var nom = prompt("Entrez votre nom : ");
        
        ws.onopen = () => {
            console.log("Connected")
            document.getElementById("sendtextchat").addEventListener("click", (e)=> {
                 e.preventDefault();
                 var data = {
                    "message" : document.getElementById("textchat").value,
                    "name" : nom
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

        function GetCookie(name) {
            var nameEQ = name + "=";
            var ca = document.cookie.split(';');
            for(var i=0;i < ca.length;i++) {
                var c = ca[i];
                while (c.charAt(0)==' ') c = c.substring(1,c.length);
                if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
            }
            return null;
        }
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

    
            function startMusicCountdown() {
                var countdownElement = document.getElementById('countdown');
                var musicCountdown = 30;
                countdownElement.innerText = musicCountdown;

                var musicCountdownInterval = setInterval(function() {
                    musicCountdown--;
                    countdownElement.innerText = musicCountdown;

                    if (musicCountdown <= 0) {
                        clearInterval(musicCountdownInterval);
                        conn.send('end_music');
                    }
                }, 1000);
            }

            conn.onmessage = function(e) {
                if (firstMessage === true) {
                    firstMessage = false;
                    return;
                }

                var data;
                try {
                    data = JSON.parse(e.data);
                } catch (error) {
                    data = e.data;
                }

                console.log(data);
                var audio = document.getElementById('audio');
                var countdownElement = document.getElementById('countdown');

                if (data === 'start_music') {
                    var countdown = 3;
                    countdownElement.innerText = countdown;

                    var countdownInterval = setInterval(function() {
                        countdown--;
                        countdownElement.innerText = countdown;

                        if (countdown <= 0) {
                            clearInterval(countdownInterval);
                            audio.play();
                            startMusicCountdown();
                        }
                    }, 1000);
                } else if (data.Preview) {
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
                                startMusicCountdown();
                            }
                        }, 1000);
                    }
                }
            };
    
            document.querySelector('form').addEventListener('submit', function(e) {
                e.preventDefault();
                if (conn.readyState === WebSocket.OPEN) {
                    var answer = document.querySelector('input[name="blindtest_answer"]').value;
                    if (typeof answer === 'string') {
                        conn.send(answer);
                    } else {
                        console.error('Answer is not a string:', answer);
                    }
                } else {
                    console.error('WebSocket is not open:', conn.readyState);
                }
            });
        };