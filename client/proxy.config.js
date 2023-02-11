const PROXY_CONFIG = {
  "/api": {
      "target": "http://test.ryanwilk.com",
      "changeOrigin": true,
      "secure": false,
      "logLevel": "debug",
      "onProxyReq": function(pr, req, res) {
          pr.removeHeader('Origin');
      }
 }
};
module.exports = PROXY_CONFIG;