import Document from 'next/document'
import { ServerStyleSheet as StyledComponentSheet } from 'styled-components'
import { ServerStyleSheets as MaterialUiStyleSheets } from '@material-ui/styles'

class MyDocument extends Document {
  static async getInitialProps (ctx) {
    const styledComponentSheet = new StyledComponentSheet()
    const materialUiStyleSheets = new MaterialUiStyleSheets()
    const originalRenderPage = ctx.renderPage

    try {
      ctx.renderPage = () =>
        originalRenderPage({
          enhanceApp: App => props =>
            styledComponentSheet.collectStyles(
              materialUiStyleSheets.collect(<App {...props} />)
            )
        })

      const initialProps = await Document.getInitialProps(ctx)

      return {
        ...initialProps,
        styles: [
          initialProps.styles,
          styledComponentSheet.getStyleElement(),
          materialUiStyleSheets.getStyleElement()
        ]
      }
    } catch (err) {
      console.error(err.stack)
    } finally {
      styledComponentSheet.seal()
    }
  }
}

export default MyDocument
