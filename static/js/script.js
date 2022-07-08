var tmp = 0;

var lt = /</g, 
    gt = />/g, 
    ap = /'/g, 
    ic = /"/g;

function repeat() {
  fetch("https://go-chatters.herokuapp.com/u/msglist").then(res => res.json()).then(d => {
    if (d.length > tmp) {
      tmp = d.length;
      document.getElementById("globalmsgs").innerText = "";
      for (let i = 0; i < d.length; i++) {
        d[i].msg = d[i].msg.toString().replace(lt, "&lt;").replace(gt, "&gt;").replace(ap, "&#39;").replace(ic, "&#34;");
        d[i].user = d[i].user.toString().replace(lt, "&lt;").replace(gt, "&gt;").replace(ap, "&#39;").replace(ic, "&#34;");
        d[i].time = d[i].time.toString().replace(lt, "&lt;").replace(gt, "&gt;").replace(ap, "&#39;").replace(ic, "&#34;");
        
        document.getElementById("globalmsgs").innerHTML += "<div class='msg'><b>" + d[i].time + d[i].user + "</b>" + d[i].msg + "</div>";
      }
    }
  });
  fetch("https://go-chatters.herokuapp.com/u/userlist").then(res => res.json()).then(d => {
    if (d.length > tmp) {
      tmp = d.length;
      document.getElementById("userlist").innerText = "";
      for (let i = 0; i < d.length; i++){
        d[i].user = d[i].user.toString().replace(lt, "&lt;").replace(gt, "&gt;").replace(ap, "&#39;").replace(ic, "&#34;");
        document.getElementById("userlist").innerHTML += "<div class='user' username=d[i].user><b>"  + d[i].user + "</b>" + "</div>";
      }
    }
  });
}
repeat()
// tmpinterval = setInterval(repeat, 3000);

var b=document.getElementsByClassName("user");
for(let i=0;i<b.length;i++){
  (function (i){
    b[i].onclick=function (){
      window.location.href="/u/privatechat?touser="+this.id
      // window.open("/u/privatechat?touser="+this.id)
    };
  })(i)
}
