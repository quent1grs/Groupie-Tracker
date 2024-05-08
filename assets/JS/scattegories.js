 var ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
        console.log("Connected");
        var sendButton = document.getElementById("sendtextchat");
        var textInput = document.getElementById("textchat");

        sendButton.addEventListener("click", (e) => {
            e.preventDefault();
            sendMessage();
        });

        textInput.addEventListener("keypress", (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                sendMessage();
            }
        });

        let letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split(''); // Définissez la variable letters
        let letter = letters.splice(Math.floor(Math.random() * letters.length), 1)[0]; // Choisissez une lettre au hasard et supprimez-la du tableau
        document.getElementById('letter').textContent = letter; 
    };

    function sendMessage() {
        var data = {
            "message": document.getElementById("textchat").value,
        };
        ws.send(JSON.stringify(data));
        document.getElementById("textchat").value = ""; // Effacez le contenu de l'input
    }
    
        ws.onmessage = (e) => {
            document.getElementById("chat").innerHTML += e.data + "<br>";
        };
    
        ws.onerror = (e) => {
            console.log(e);
        };
    
        ws.onclose = () => {
            delete ws;
        };
    
        // function GetCookie(name) {
        //     var nameEQ = name + "=";
        //     var ca = document.cookie.split(';');
        //     for (var i = 0; i < ca.length; i++) {
        //         var c = ca[i];
        //         while (c.charAt(0) == ' ') c = c.substring(1, c.length);
        //         if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
        //     }
        //     return null;
        // }
        
        $(document).ready(function() {
            $.ajax({
                url: '/getLetter',
                type: 'GET',
                success: function(data) {
                    $('#letter').text(data.letter);
                },
                error: function(error) {
                    console.log(error);
                }
            });
        });