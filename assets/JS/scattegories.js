 var ws = new WebSocket('ws://localhost:8080/chatscattegoriesws');
 var ScattegoriesGameSocket = new WebSocket('ws://localhost:8080/ScattegoriesGame');

        ws.onopen = () => {
        console.log("Chat Connected");
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
            // "username": ,
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
        document.addEventListener("DOMContentLoaded", function() {
            var ScattegoriesGameSocket = new WebSocket('ws://localhost:8080/ScattegoriesGame');
        
            ScattegoriesGameSocket.onopen = function(event) {
                console.log("ScattegoriesGameSocket is open");
            };
        
            ScattegoriesGameSocket.onclose = function(event) {
                console.log("ScattegoriesGameSocket is closed");
            };
        
            ScattegoriesGameSocket.onerror = function(error) {
                console.error("Error on ScattegoriesGameSocket: ", error);
            };
        
            // Le reste de votre code...
        });
