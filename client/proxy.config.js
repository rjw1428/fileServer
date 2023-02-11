const PROXY_CONFIG = {
  "/api": {
    //   "target": "http://test.ryanwilk.com",
    "target": "http://localhost:1428",
      "changeOrigin": true,
      "secure": false,
      "logLevel": "debug",
      "onProxyReq": function(pr, req, res) {
          pr.removeHeader('Origin');
      }
 }
};
module.exports = PROXY_CONFIG;