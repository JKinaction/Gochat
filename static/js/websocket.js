
var conn;
var log=document.getElementById("log");
$(document).ready(function (){

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    if (window["WebSocket"]){
        conn=new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage=function (evt){
            var parse = JSON.parse(evt.data);
            console.log(parse)
            if (parse.status==1||parse.status==2){
                var message=parse.data.content.split('\n');
                var time=parse.data.time.split('\n');
                var username=parse.data.user.split('\n');
                for (var i=0;i<message.length;i++){
                    document.getElementById("globalmsgs").innerHTML += "<div class='msg'><b>" + time[i]+ username[i] +":" + "</b>" + message[i] + "</div>";
                }
            }
        }
    }else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
})