module.exports = {
  env: {
    API_SERVER_NAME: process.env.NODE_ENV === 'production'
      ? `http://${process.env.API_SERVER_HOST}:${process.env.API_SERVER_PORT}`
      : 'http://localhost:3010'
  },
}
