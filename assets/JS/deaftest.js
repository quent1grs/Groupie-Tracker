window.onload = function() {
    var conn = new WebSocket('ws://localhost:8080/deaftestws');
    console.log(conn);

    var intervalId;
    var timer = 30;

    conn.onopen = function(e) {
        console.log("Connection established!");
    };

    conn.onclose = function(e) {
        console.log("Connection closed!");
    };

    conn.onmessage = function(e) {
        var data;
        try {
            data = JSON.parse(e.data);
        } catch (error) {
            data = e.data;
        }

        var lyrics = document.getElementById('Lyrics');
        var countdown = document.getElementById('countdown');
        var previous = document.getElementById('previous');

        if (data.Lyrics) {
            var formattedLyrics = data.Lyrics.replace(/\n/g, '<br>');
            lyrics.innerHTML = formattedLyrics;
        } else if (data.Title && data.Artist) {
            lyrics.innerText = "Title: " + data.Title + " by " + data.Artist;
        }

        timer = 30;
        if (intervalId) {
            clearInterval(intervalId);
        }

        intervalId = setInterval(function() {
            timer--;
            countdown.innerText = timer;
            if (timer === 0) {
                clearInterval(intervalId);
                previous.innerText = "Previous song: " + data.Title + " by " + data.Artist;
                conn.send('Change_song');
            }
        }, 1000);
    };

    document.querySelector('form').addEventListener('submit', function(e) {
        e.preventDefault();
        if (conn.readyState === WebSocket.OPEN) {
            remainingTime = parseInt(document.getElementById('countdown').innerText, 10);
            var answer = document.querySelector('input[name="deaftest_answer"]').value;
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