const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: "turn:relay.metered.ca:443",
      username: 'a9067dff0bdee1097e961805',
      credential: 'btIRqKUbbhxsazf3'
    }
  ],
      iceCandidatePoolSize: 10,
      sdpSemantics: 'unified-plan' //newer implementation of WebRTC
})

let localDescription = null

pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)
pc.onicecandidate = event => {
  if (event.candidate === null) { // && event.candidate.candidate.includes("TCP")) {
console.log(event)
    localDescription = event.target.localDescription
  }
}

pc.addTransceiver('audio')
pc.createOffer()
  .then(d => pc.setLocalDescription(d))
  .catch(console.log)

pc.ontrack = function (event) {
  const el = document.getElementById('audio1')
  el.srcObject = event.streams[0]
  console.log(pc);
  el.autoplay = true
  el.controls = true
  el.value = 1;
  el.muted = false
}

window.startSession = () => {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "http://5.159.101.107:8080/channel/2/connect");
  xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');
  console.log(localDescription);
  xhr.send(JSON.stringify({sdp: localDescription.sdp}));
  xhr.onload = (e) => {
    response = JSON.parse(xhr.response)
    console.log(response)
    try {
      pc.setRemoteDescription(response)
    } catch (e) {
      alert(e)
    }
  }
}
