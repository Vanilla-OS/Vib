const express = require("express");
const path = require("path");
const app = express();
const port = 3000;

app.use(express.static("dist"));

app.use((req, res, next) => {
  if (req.path.indexOf(".") === -1) {
    const file = path.join(__dirname, "dist", `${req.path}.html`);
    res.sendFile(file);
  } else {
    next();
  }
});

app.listen(port, () => {
  console.log(`Proxy server listening at http://localhost:${port}`);
});
