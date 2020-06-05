module.exports = {
  env: {
    API_SERVER_NAME: process.env.NODE_ENV === 'production'
      ? `http://${process.env.API_SERVER_HOST || 'localhost'}:${process.env.API_SERVER_PORT || 3010}`
      : 'http://localhost:3010'
  },
}
