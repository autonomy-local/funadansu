// 旧 operator の 07-01 が jest 内で立てる catalog のスタブを、別プロセスで再現する（pxr-setting.json を返す）
const http = require('http');
const fs = require('fs');
const setting = fs.readFileSync(process.argv[2], 'utf-8');
http.createServer((req, res) => {
  const url = new URL(req.url, 'http://x');
  if (url.pathname === '/catalog' && url.searchParams.get('ns') === 'catalog/ext/test-org/setting/global') {
    res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(setting);
  } else { res.writeHead(204); res.end(); }
}).listen(3001, () => console.log('catalog stub on 3001'));
