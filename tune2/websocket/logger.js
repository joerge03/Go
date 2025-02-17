if (window["WebSocket"]){
    ws = new WebSocket("ws://{{.}}/ws")
    // console.log(e.key)
    document.onkeydown = (e)=>{
        console.log(e.key)
        ws.send(e.key)
    }
}
