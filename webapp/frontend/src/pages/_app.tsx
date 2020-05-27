import { useEffect } from 'react'
import Head from 'next/head'
import { ThemeProvider, makeStyles, createStyles } from '@material-ui/core/styles'
import CssBaseline from '@material-ui/core/CssBaseline'
import Paper from '@material-ui/core/Paper'
import theme from '../plugins/theme'

import type { FC } from 'react'
import type { AppProps } from 'next/app'
import type { Theme } from '@material-ui/core/styles'

import 'leaflet/dist/leaflet.css'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    page: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    }
  })
)

const MyApp: FC<AppProps> = props => {
  const { Component, pageProps } = props

  useEffect(() => {
    // Remove the server-side injected CSS.
    const jssStyles = document.querySelector('#jss-server-side')
    if (jssStyles) {
      jssStyles.parentElement?.removeChild(jssStyles)
    }

    import('leaflet')
      .then(Leaflet => {
        delete (Leaflet.Icon.Default.prototype as any)._getIconUrl
        Leaflet.Icon.Default.mergeOptions({
          iconRetinaUrl: '/images/leaflet/marker-icon-2x.png',
          iconUrl: '/images/leaflet/marker-icon.png',
          shadowUrl: '/images/leaflet/marker-shadow.png'
        })
      })
      .catch(error => console.error('failed to set marker icons.', error))
  }, [])

  const Page: FC = props => {
    const classes = useStyles()

    return (
      <Paper className={classes.page}>
        <Component {...props} />
      </Paper>
    )
  }

  return (
    <>
      <Head>
        <title>isuumo</title>
        <meta name='viewport' content='minimum-scale=1, initial-scale=1, width=device-width' />
      </Head>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Page {...pageProps} />
      </ThemeProvider>
    </>
  )
}

export default MyApp
