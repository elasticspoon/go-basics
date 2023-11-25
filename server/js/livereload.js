let socket = new WebSocket('ws://localhost:8080/alive');

function reloadPage() {
  console.log('Reloading page');
  window.location.reload();
  socket.removeEventListener('open', reloadPage);
}

function checkConnection() {
  socket = new WebSocket('ws://localhost:8080/alive');
  socket.addEventListener('open', reloadPage);

  socket.addEventListener('error', () => {
    setTimeout(checkConnection, 1000);
  });
}

socket.addEventListener('close', () => {
  checkConnection();
});

