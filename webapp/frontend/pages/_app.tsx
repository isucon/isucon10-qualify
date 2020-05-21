import App from 'next/app'
import { ThemeProvider } from '@material-ui/styles'
import { CssBaseline } from '@material-ui/core'

import Mytheme from '../plugins/theme'

class MyApp extends App {
  componentDidMount () {
    const jssStyles = document.querySelector('#jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  render () {
    const { Component, pageProps } = this.props

    return (
      <ThemeProvider theme={Mytheme}>
        <CssBaseline />
        <Component {...pageProps} />
      </ThemeProvider>
    )
  }
}

export default MyApp
